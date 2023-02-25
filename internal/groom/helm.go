// package groom is used to add or remove poker room
// in the context of helm deployment with .Values.rooms
package groom

import (
	"context"
	"errors"
	"log"
	"os"

	zlog "github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

const (
	releaseName    = "poker"
	roomsValueName = "rooms"
)

// get list of deployed rooms
func RoomDeployed() (roomDeployed []string, err error) {
	pokerRelease, err := getPokerRelease()
	if err != nil {
		zlog.Error().Msg("unable to get helm release for poker")
		return []string{}, err
	}

	rooms := []string{}
	// get rooms
	for _, room := range pokerRelease.Config[roomsValueName].([]interface{}) {
		rooms = append(rooms, room.(string))
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
		if room == roomName {
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
		zlog.Error().Msg("unable to get helm releases")
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
	return updateRoom(roomName, true)
}

// remove room named roomName
func RemoveRoom(roomName string) (rooms []string, err error) {
	// check if room already exists
	exists, err := roomExists(roomName)
	if err != nil {
		return nil, err
	}
	if exists {
		return updateRoom(roomName, false)
	}
	return rooms, errors.New("Room named " + roomName + " not found. No change applied.")
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
		zlog.Error().Msg("unable to init helm client for upgrade")
		return nil, err
	}
	if err = cfg.KubeClient.IsReachable(); err != nil {
		zlog.Error().Msg("unable to connect to the kubernetes cluster with helm client for upgrade")
		return nil, err
	}
	client := action.NewUpgrade(cfg)

	client.Namespace = settings.Namespace()
	if err := chartutil.ValidateReleaseName(releaseName); err != nil {
		zlog.Error().Msg("unable to validate release name with helm client for upgrade")
		return nil, err
	}
	ctx := context.Background()

	newValues := prepareValues(pokerRelease, isAdd, roomName)
	zlog.Info().Msgf("Values target is: %v", newValues)

	// apply changes on kubernetes namespace
	if _, err := client.RunWithContext(ctx, releaseName, pokerRelease.Chart, newValues); err != nil {
		zlog.Error().Msg("helm client for upgrade FAILED to apply release update")
		return nil, err
	}
	for _, room := range newValues[roomsValueName].([]interface{}) {
		rooms = append(rooms, room.(string))
	}
	return rooms, nil
}

// prepare values to apply to the current helm release to add/remove one room
func prepareValues(pokerRelease *release.Release, isAdd bool, roomName string) map[string]interface{} {
	newValues := make(map[string]interface{})
	for k, v := range pokerRelease.Config {
		newValues[k] = v
		if k == roomsValueName {
			// add roomName
			if isAdd {
				newValues[k] = append(newValues[k].([]interface{}), roomName)
			} else {
				// delete roomName
				newValues[k] = []interface{}{} // remove all rooms
				for _, name := range v.([]interface{}) {
					// re-add other rooms only
					if name != roomName {
						newValues[k] = append(newValues[k].([]interface{}), name)
					}
				}
			}
		}
	}
	return newValues
}
