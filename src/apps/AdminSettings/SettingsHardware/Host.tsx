import { Host as HostModel } from "../backend/models";
import { SiApple, SiLinux, SiWindows } from "@icons-pack/react-simple-icons";
import {
    Button,
    List,
    ListActions,
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
    MaterialIcon,
    Paragraph,
} from "@vertex-center/components";
import { Fragment, useState } from "react";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../components/KeyValueInfo/KeyValueInfo";
import { useReboot } from "../hooks/useReboot";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { APIError } from "../../../components/Error/APIError";
import Popup from "../../../components/Popup/Popup";

type HostProps = {
    host?: HostModel;
};

export default function Host(props: Readonly<HostProps>) {
    if (!props.host) return null;

    const [showRebootPopup, setShowRebootPopup] = useState(false);
    const [rebootSent, setRebootSent] = useState(false);

    const { reboot, isRebooting, errorReboot } = useReboot({
        onSuccess: () => {
            setRebootSent(true);
        },
    });

    const {
        os,
        hostname,
        platform,
        platform_version,
        kernel_arch,
        uptime,
        boot_time,
    } = props.host;

    let icon = undefined;
    switch (os) {
        case "linux":
            icon = <SiLinux />;
            break;
        case "darwin":
            icon = <SiApple />;
            break;
        case "windows":
            icon = <SiWindows />;
            break;
    }

    const uptimeHours = Math.round((uptime / 3600) * 100) / 100;
    const bootTime = new Date(boot_time * 1000).toLocaleString();

    const onPopupDismiss = () => setShowRebootPopup(false);

    const popupActions = (
        <Button
            variant="danger"
            leftIcon={<MaterialIcon icon="restart_alt" />}
            onClick={() => reboot()}
            disabled={rebootSent || isRebooting}
        >
            Reboot
        </Button>
    );

    return (
        <Fragment>
            <ProgressOverlay show={isRebooting} />
            <APIError error={errorReboot} />
            <List>
                <ListItem>
                    <ListIcon>{icon}</ListIcon>
                    <ListInfo>
                        <ListTitle>{hostname}</ListTitle>
                        <ListDescription>
                            {platform}({platform_version}) - {kernel_arch}
                        </ListDescription>
                    </ListInfo>
                    <ListActions>
                        <Button
                            variant="danger"
                            leftIcon={<MaterialIcon icon="restart_alt" />}
                            onClick={() => setShowRebootPopup(true)}
                            disabled={
                                rebootSent || isRebooting || showRebootPopup
                            }
                        >
                            Reboot
                        </Button>
                    </ListActions>
                </ListItem>
            </List>
            <KeyValueGroup>
                <KeyValueInfo name="Boot Time" icon="restart_alt">
                    {bootTime}
                </KeyValueInfo>
                <KeyValueInfo name="Uptime" icon="arrow_upward">
                    {uptimeHours}h
                </KeyValueInfo>
            </KeyValueGroup>
            <Popup
                show={showRebootPopup}
                onDismiss={onPopupDismiss}
                title="Reboot"
                actions={popupActions}
            >
                <Paragraph>
                    Are you sure you want to reboot <strong>{hostname}</strong>?
                    This will disconnect you from the server.
                </Paragraph>
            </Popup>
        </Fragment>
    );
}
