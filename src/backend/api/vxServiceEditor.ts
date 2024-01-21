import { Template } from "../../apps/Containers/backend/template";

import { createServer } from "../server";

// @ts-ignore
const server = createServer(window.api_urls.devtools_service_editor);

const toYaml = async (service: Template) => {
    const { data } = await server.post(`/editor/to-yaml`, service);
    return data;
};

const editorRoutes = {
    toYaml: toYaml,
};

export const vxServiceEditorRoutes = {
    editor: editorRoutes,
};
