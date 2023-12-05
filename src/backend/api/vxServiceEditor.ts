import { Service } from "../../models/service";

import { createServer } from "../server";

const server = createServer("7510");

const toYaml = async (service: Service) => {
    const { data } = await server.post(`/editor/to-yaml`, service);
    return data;
};

const editorRoutes = {
    toYaml: toYaml,
};

export const vxServiceEditorRoutes = {
    editor: editorRoutes,
};
