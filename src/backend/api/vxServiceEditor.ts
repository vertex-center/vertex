import { Service } from "../../apps/Containers/backend/service";

import { createServer } from "../server";

// @ts-ignore
const server = createServer(window.api_urls.devtools_service_editor);

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
