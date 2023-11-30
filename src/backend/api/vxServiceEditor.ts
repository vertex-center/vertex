import { Service } from "../../models/service";
import { server } from "./backend";

const BASE_URL = `/app/devtools-service-editor`;

const toYaml = async (service: Service) => {
    const { data } = await server.post(`${BASE_URL}/editor/to-yaml`, service);
    return data;
};

const editorRoutes = {
    toYaml: toYaml,
};

export const vxServiceEditorRoutes = {
    editor: editorRoutes,
};
