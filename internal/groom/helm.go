// package groom is used to add or remove poker room
// in the context of helm deployment with .Values.rooms
package groom

import (
	"context"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"

	"github.com/fc92/poker/internal/common/logger"
)

const (
	releaseName    = "poker"
	roomsValueName = "rooms"
)

var index int

func init() {
	logger.InitLogger()
	index = 30
}

// get list of deployed rooms
func RoomDeployed() (roomDeployed []interface{}, err error) {
	pokerRelease, err := getPokerRelease()
	if err != nil {
		log.Logger.Error().Msg("unable to get helm release for poker")
		return []interface{}{}, err
	}

	rooms := []interface{}{}
	// get rooms
	for _, room := range pokerRelease.Config[roomsValueName].([]interface{}) {
		rooms = append(rooms, room)
	}
	return rooms, nil
}

// check if room named roomName exists
func roomExists(roomName string) (exists bool, err error) {
	roomDeployed, err := RoomDeployed()
	if err != nil {
		return false, err
	}
	for _, room := range roomDeployed {
		roomMap := room.(map[string]interface{})
		if roomMap["name"] == roomName {
			return true, nil
		}
	}
	return false, nil
}

// get the release used for poker deployment
func getPokerRelease() (pokerRelease *release.Release, err error) {
	// get list of Helm chart releases deployed
	releases, err := getReleases()
	if err != nil {
		log.Logger.Error().Msg("unable to get helm releases")
		return nil, err
	}

	// identify release for poker
	for _, rel := range releases {
		if rel.Name == releaseName {
			return rel, nil
		}
	}
	return nil, errors.New(releaseName + " release not found")
}

// get helm releases deployed on the namespace
func getReleases() (releases []*release.Release, err error) {
	settings := cli.New()

	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return []*release.Release{}, err
	}

	clientRead := action.NewList(actionConfig)

	clientRead.Deployed = true
	results, err := clientRead.Run()
	if err != nil {
		return []*release.Release{}, err
	}
	return results, nil
}

// add room named roomName
func AddRoom(roomName string) (rooms []string, err error) {
	// check if room already exists
	exists, err := roomExists(roomName)
	if err != nil {
		return nil, err
	}
	if exists {
		return rooms, errors.New("Room named " + roomName + " already exists. No change applied.")
	}
	log.Logger.Info().Msgf("Adding room %v", roomName)
	return updateRoom(roomName, true)
}

// remove room named roomName
func RemoveRoom(roomName string) {
	// check if room already exists
	exists, err := roomExists(roomName)
	if err != nil {
		log.Logger.Error().Msgf("checking room %s exists failed with error %v", roomName, err)
	}
	if exists {
		log.Logger.Info().Msgf("Removing room %v", roomName)
		_, err = updateRoom(roomName, false)
		if err != nil {
			log.Logger.Error().Msgf("removal of room %s failed with error %v", roomName, err)
		}
	}
	log.Logger.Error().Msgf("removal of room %s failed", roomName)
}

// add or remove room based on isAdd value
func updateRoom(roomName string, isAdd bool) (rooms []string, err error) {
	pokerRelease, err := getPokerRelease()
	if err != nil {
		return nil, err
	}
	settings := cli.New()
	cfg := new(action.Configuration)
	if err := cfg.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Logger.Error().Msg("unable to init helm client for upgrade")
		return nil, err
	}
	if err = cfg.KubeClient.IsReachable(); err != nil {
		log.Logger.Error().Msg("unable to connect to the kubernetes cluster with helm client for upgrade")
		return nil, err
	}
	client := action.NewUpgrade(cfg)
	client.ReuseValues = true

	client.Namespace = settings.Namespace()
	if err := chartutil.ValidateReleaseName(releaseName); err != nil {
		log.Logger.Error().Msg("unable to validate release name with helm client for upgrade")
		return nil, err
	}
	ctx := context.Background()

	newValues := prepareValues(pokerRelease, isAdd, roomName)
	log.Logger.Info().Msgf("Values target is: %v", newValues)

	// apply changes on kubernetes namespace
	if _, err := client.RunWithContext(ctx, releaseName, pokerRelease.Chart, newValues); err != nil {
		log.Logger.Error().Msg("helm client for upgrade FAILED to apply release update")
		return nil, err
	}
	for _, room := range newValues[roomsValueName].([]interface{}) {
		roomMap := room.(map[string]interface{})
		rooms = append(rooms, roomMap["name"].(string))
	}
	return rooms, nil
}

func prepareValues(pokerRelease *release.Release, isAdd bool, roomName string) map[string]interface{} {
	newValues := make(map[string]interface{})
	for k, v := range pokerRelease.Config {
		newValues[k] = v
		if k == roomsValueName {
			// add roomName
			if isAdd {
				newRoom := make(map[string]interface{})
				index++
				newRoom["name"] = roomName
				newRoom["index"] = index

				newValues[k] = append(newValues[k].([]interface{}), newRoom)
			} else {
				// delete roomName
				newValues[k] = []interface{}{} // remove all rooms
				for _, room := range v.([]interface{}) {
					roomMap := room.(map[string]interface{})
					// re-add other rooms only
					if roomMap["name"] != roomName {
						newValues[k] = append(newValues[k].([]interface{}), room)
					}
				}
			}
		}
	}
	return newValues
}
