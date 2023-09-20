export type Dependencies = {
    last_updates_check: string;
    items: Dependency[];
};

export type Dependency = {
    id: string;
    name: string;
    version: string;
    update?: DependencyUpdate;
};

export type DependencyUpdate = {
    current_version: string;
    latest_version: string;
    needs_restart?: boolean;
};
