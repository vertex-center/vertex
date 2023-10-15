import { useState } from "react";
import { Caption, Text, Title } from "../../../components/Text/Text";
import { Horizontal, Vertical } from "../../../components/Layouts/Layouts";
import Button from "../../../components/Button/Button";
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
import Icon from "../../../components/Icon/Icon";

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

    const isLoading =
        isLoadingUpdate ||
        isLoadingSettings ||
        isInstallingUpdate ||
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
            {update === undefined && !isLoadingUpdate && (
                <Caption className={styles.content}>
                    <Horizontal alignItems="center" gap={6}>
                        <Icon name="check" />
                        Vertex is up to date. You are running the latest
                        version.
                    </Horizontal>
                </Caption>
            )}
            {update !== undefined && (
                <List>
                    <VertexUpdate
                        version={update?.baseline.version}
                        description={update?.baseline.description}
                        install={installUpdate}
                        isInstalling={isInstallingUpdate}
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
                    <Button primary onClick={dismissPopup} rightIcon="check">
                        OK
                    </Button>
                </Horizontal>
            </Popup>
        </Vertical>
    );
}
