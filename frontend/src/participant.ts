export interface Participant {
    id: string;
    name: string;
    vote: string;
    available_commands: Record<string, string>;
    last_command: string;
}