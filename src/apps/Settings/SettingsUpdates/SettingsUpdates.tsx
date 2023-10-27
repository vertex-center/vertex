import { useState } from "react";
import { Caption, Text, Title } from "../../../components/Text/Text";
import { Horizontal, Vertical } from "../../../components/Layouts/Layouts";
import { Button, MaterialIcon } from "@vertex-center/components";
import Spacer from "../../../components/Spacer/Spacer";
import Popup from "../../../components/Popup/Popup";
import styles from "./SettingsUpdates.module.sass";
import VertexUpdate from "../components/VertexUpdate/VertexUpdate";
import { APIError } from "../../../components/Error/APIError";
import ToggleButton from "../../../components/ToggleButton/ToggleButton";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useQueryClient } from "@tanstack/react-query";
import List from "../../../components/List/List";
import { useSettings } from "../hooks/useSettings";
import { useUpdate } from "../hooks/useUpdate";
import { useUpdateMutation } from "../hooks/useUpdateMutation";
import { useSettingsChannelMutation } from "../hooks/useSettingsMutation";

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

    const { setChannel, isSettingChannel, errorSetChannel } =
        useSettingsChannelMutation({
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
        isSettingChannel;

    const error =
        errorUpdate || errorSettings || errorInstallUpdate || errorSetChannel;

    return (
        <Vertical gap={20}>
            <ProgressOverlay show={isLoading} />
            <Title className={styles.title}>Updates</Title>
            <Horizontal className={styles.toggle} alignItems="center">
                <Text>Enable Beta channel</Text>
                <Spacer />
                <ToggleButton
                    value={settings?.updates?.channel === "beta"}
                    onChange={(beta: boolean) => setChannel(beta)}
                    disabled={isLoading}
                />
            </Horizontal>
            <APIError error={error} />
            {update === null && !isLoadingUpdate && (
                <Caption className={styles.content}>
                    <Horizontal alignItems="center" gap={6}>
                        <MaterialIcon icon="check" />
                        Vertex is up to date. You are running the latest
                        version.
                    </Horizontal>
                </Caption>
            )}
            {update !== null && (
                <List>
                    <VertexUpdate
                        version={update?.baseline.version}
                        description={update?.baseline.description}
                        install={installUpdate}
                        isInstalling={isInstalling}
                    />
                </List>
            )}
            <Popup show={showMessage} onDismiss={dismissPopup}>
                <Text>
                    Updates are installed. You can now restart your Vertex
                    server.
                </Text>
                <Horizontal>
                    <Spacer />
                    <Button
                        variant="colored"
                        onClick={dismissPopup}
                        rightIcon={<MaterialIcon icon="check" />}
                    >
                        OK
                    </Button>
                </Horizontal>
            </Popup>
        </Vertical>
    );
}
