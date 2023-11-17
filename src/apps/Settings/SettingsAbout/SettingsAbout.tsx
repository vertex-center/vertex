import { api } from "../../../backend/api/backend";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../components/KeyValueInfo/KeyValueInfo";
import { APIError } from "../../../components/Error/APIError";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useQuery } from "@tanstack/react-query";
import { Title } from "@vertex-center/components";
import Content from "../../../components/Content/Content";

export default function SettingsAbout() {
    const {
        data: about,
        isLoading,
        error,
    } = useQuery({
        queryKey: ["about"],
        queryFn: api.about,
    });

    return (
        <Content>
            <Title variant="h2">Vertex</Title>
            <ProgressOverlay show={isLoading} />
            <APIError error={error} />
            <KeyValueGroup>
                <KeyValueInfo
                    name="Version"
                    type="code"
                    icon="tag"
                    loading={isLoading}
                >
                    {about?.version}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Commit"
                    type="code"
                    icon="commit"
                    loading={isLoading}
                >
                    {about?.commit}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Release date"
                    type="code"
                    icon="calendar_month"
                    loading={isLoading}
                >
                    {about?.date}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Compiled for"
                    type="code"
                    icon="memory"
                    loading={isLoading}
                >
                    {about?.os}
                    {about?.arch && `/${about?.arch}`}
                </KeyValueInfo>
            </KeyValueGroup>
        </Content>
    );
}
