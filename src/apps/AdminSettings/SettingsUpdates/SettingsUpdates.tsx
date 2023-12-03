import { useState } from "react";
import { Caption } from "../../../components/Text/Text";
import { Horizontal } from "../../../components/Layouts/Layouts";
import {
    Button,
    List,
    MaterialIcon,
    Paragraph,
    Title,
} from "@vertex-center/components";
import Spacer from "../../../components/Spacer/Spacer";
import Popup from "../../../components/Popup/Popup";
import VertexUpdate from "../components/VertexUpdate/VertexUpdate";
import { APIError } from "../../../components/Error/APIError";
import ToggleButton from "../../../components/ToggleButton/ToggleButton";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useQueryClient } from "@tanstack/react-query";
import { useSettings } from "../hooks/useSettings";
import { useUpdate } from "../hooks/useUpdate";
import { useUpdateMutation } from "../hooks/useUpdateMutation";
import { usePatchSettings } from "../hooks/usePatchSettings";
import Content from "../../../components/Content/Content";

export default function SettingsUpdates() {
    const queryClient = useQueryClient();
    const [showMessage, setShowMessage] = useState<boolean>(false);

    const { update, isLoadingUpdate, errorUpdate } = useUpdate();
    const { settings, isLoadingSettings, errorSettings } = useSettings();

    const { installUpdate, isInstallingUpdate, errorInstallUpdate } =
        useUpdateMutation({
            onSuccess: () => {
                setShowMessage(true);
                queryClient.invalidateQueries({
                    queryKey: ["updates"],
                });
            },
        });

    const { patchSettings, isPatchingSettings, errorPatchingSettings } =
        usePatchSettings({
            onSuccess: () => {
                queryClient.invalidateQueries({
                    queryKey: ["settings"],
                });
            },
        });

    const dismissPopup = () => {
        setShowMessage(false);
    };

    const isInstalling = update?.updating === true || isInstallingUpdate;

    const isLoading =
        isLoadingUpdate ||
        isLoadingSettings ||
        isInstalling ||
        isPatchingSettings;

    const error =
        errorUpdate ||
        errorSettings ||
        errorInstallUpdate ||
        errorPatchingSettings;

    const actions = (
        <Button
            variant="colored"
            onClick={dismissPopup}
            rightIcon={<MaterialIcon icon="check" />}
        >
            OK
        </Button>
    );

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
                <List>
                    <VertexUpdate
                        version={update?.baseline.version}
                        description={update?.baseline.description}
                        install={installUpdate}
                        isInstalling={isInstalling}
                    />
                </List>
            )}
            <Popup
                show={showMessage}
                onDismiss={dismissPopup}
                title="Updates are installed"
                actions={actions}
            >
                <Paragraph>You can now restart your Vertex server.</Paragraph>
            </Popup>
        </Content>
    );
}
