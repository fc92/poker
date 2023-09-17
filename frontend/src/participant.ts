import { v4 as uuidv4 } from 'uuid';

export interface Participant {
    id: string; // Vous pouvez également utiliser le type UUID si vous avez une bibliothèque compatible.
    name: string;
    vote: string;
    available_commands: Record<string, string>; // Utilisation de Record pour définir un objet avec des clés de chaînes.
    last_command: string;
}