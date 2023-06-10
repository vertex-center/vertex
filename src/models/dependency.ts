import { Package } from "./package";

export type Dependency = Package & {
    installed: boolean;
};

export type Dependencies = { [id: string]: Dependency };
