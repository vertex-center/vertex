import { Caption } from "../../../components/Text/Text";
import { Horizontal } from "../../../components/Layouts/Layouts";
import {
    List,
    MaterialIcon,
    Paragraph,
    Title,
    Vertical,
} from "@vertex-center/components";
import Spacer from "../../../components/Spacer/Spacer";
import VertexUpdate from "../components/VertexUpdate/VertexUpdate";
import { APIError } from "../../../components/Error/APIError";
import ToggleButton from "../../../components/ToggleButton/ToggleButton";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useQueryClient } from "@tanstack/react-query";
import { useSettings } from "../hooks/useSettings";
import { useUpdate } from "../hooks/useUpdate";
import { usePatchSettings } from "../hooks/usePatchSettings";
import Content from "../../../components/Content/Content";

export default function SettingsUpdates() {
    const queryClient = useQueryClient();

    const { update, isLoadingUpdate, errorUpdate } = useUpdate();
    const { settings, isLoadingSettings, errorSettings } = useSettings();

    const { patchSettings, isPatchingSettings, errorPatchingSettings } =
        usePatchSettings({
            onSuccess: () => {
                queryClient.invalidateQueries({
                    queryKey: ["settings"],
                });
            },
        });

    const isLoading =
        isLoadingUpdate || isLoadingSettings || isPatchingSettings;

    const error = errorUpdate || errorSettings || errorPatchingSettings;

    const hasUpdate = update !== null && update !== undefined;

    const onChannelChange = (beta: boolean) => {
        patchSettings(
            beta ? { updates_channel: "beta" } : { updates_channel: "stable" }
        );
    };

    return (
        <Content>
            <ProgressOverlay show={isLoading} />
            <Title variant="h2">Updates</Title>
            <Horizontal alignItems="center">
                <Paragraph>Enable Beta channel</Paragraph>
                <Spacer />
                <ToggleButton
                    value={settings?.updates_channel === "beta"}
                    onChange={onChannelChange}
                    disabled={isLoading}
                />
            </Horizontal>
            <APIError error={error} />
            {!hasUpdate && !isLoadingUpdate && (
                <Caption>
                    <Horizontal alignItems="center" gap={6}>
                        <MaterialIcon icon="check" />
                        Vertex is up to date. You are running the latest
                        version.
                    </Horizontal>
                </Caption>
            )}
            {hasUpdate && (
                <Vertical gap={20}>
                    <Paragraph>
                        A new version of Vertex is available. Update your
                        containers to get the latest features and bug fixes.
                    </Paragraph>
                    <List>
                        <VertexUpdate
                            version={update?.baseline.version}
                            description={update?.baseline.description}
                        />
                    </List>
                </Vertical>
            )}
        </Content>
    );
}
