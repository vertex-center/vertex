export type Update = {
    id: string;
    name: string;
    current_version: string;
    latest_version: string;
    needs_restart?: boolean;
};

export type Updates = {
    last_checked: string;
    items: Update[];
};
