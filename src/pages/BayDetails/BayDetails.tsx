import Bay from "../../components/Bay/Bay";
import { useEffect, useState } from "react";
import { api } from "../../backend/backend";
import { Outlet, useNavigate, useParams } from "react-router-dom";

import styles from "./BayDetails.module.sass";
import { Horizontal } from "../../components/Layouts/Layouts";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";
import Spacer from "../../components/Spacer/Spacer";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../components/Sidebar/Sidebar";
import Popup from "../../components/Popup/Popup";
import { Text, Title } from "../../components/Text/Text";
import Button from "../../components/Button/Button";
import Progress from "../../components/Progress";
import { SiDocker } from "@icons-pack/react-simple-icons";
import useInstance from "../../hooks/useInstance";
import { Error } from "../../components/Error/Error";

export default function BayDetails() {
    const { uuid } = useParams();
    const navigate = useNavigate();

    const { instance, setInstance } = useInstance(uuid);

    const [showDeletePopup, setShowDeletePopup] = useState<boolean>();
    const [deleting, setDeleting] = useState<boolean>(false);
    const [error, setError] = useState<string>();

    useEffect(() => {
        if (uuid === undefined) return;

        const sse = registerSSE(`/instance/${uuid}/events`);

        const onStatusChange = (e) => {
            setInstance((instance) => ({ ...instance, status: e.data }));
        };

        registerSSEEvent(sse, "status_change", onStatusChange);

        return () => {
            unregisterSSEEvent(sse, "status_change", onStatusChange);

            unregisterSSE(sse);
        };
    }, [uuid]);

    const toggleInstance = async (uuid: string) => {
        if (instance.status === "off" || instance.status === "error") {
            await api.instance.start(uuid);
        } else {
            await api.instance.stop(uuid);
        }
    };

    const onDeleteInstance = () => {
        setDeleting(true);
        setError(undefined);
        api.instance
            .delete(uuid)
            .then(() => {
                navigate("/infrastructure");
            })
            .catch((err) => {
                setError(err?.response?.data?.message ?? err?.message);
                setDeleting(false);
            });
    };

    const dismissDeletePopup = () => {
        setShowDeletePopup(false);
    };

    return (
        <div className={styles.details}>
            <div className={styles.bay}>
                <Bay
                    instances={[
                        {
                            name:
                                instance?.display_name ??
                                instance?.service?.name,
                            status: instance?.status,
                            onPower: () => toggleInstance(uuid),
                            method: instance?.install_method,
                        },
                    ]}
                />
            </div>
            <Horizontal className={styles.content}>
                <Sidebar root={`/infrastructure/${uuid}`}>
                    <SidebarGroup>
                        <SidebarItem
                            to={`/infrastructure/${uuid}/home`}
                            symbol="home"
                            name="Home"
                        />
                    </SidebarGroup>
                    <SidebarGroup title="Analyze">
                        <SidebarItem
                            to={`/infrastructure/${uuid}/logs`}
                            symbol="terminal"
                            name="Logs"
                        />
                        {instance?.install_method === "docker" && (
                            <SidebarItem
                                to={`/infrastructure/${uuid}/docker`}
                                symbol={<SiDocker size={20} />}
                                name="Docker"
                            />
                        )}
                    </SidebarGroup>
                    <SidebarGroup title="Manage">
                        <SidebarItem
                            to={`/infrastructure/${uuid}/environment`}
                            symbol="tune"
                            name="Environment"
                        />
                        {instance?.install_method !== "docker" && (
                            <SidebarItem
                                to={`/infrastructure/${uuid}/dependencies`}
                                symbol="widgets"
                                name="Dependencies"
                            />
                        )}
                        {instance?.service?.databases && (
                            <SidebarItem
                                to={`/infrastructure/${uuid}/database`}
                                symbol="database"
                                name="Database"
                            />
                        )}
                        <SidebarItem
                            to={`/infrastructure/${uuid}/update`}
                            symbol="update"
                            name="Update"
                            notifications={
                                instance?.service_update?.available
                                    ? 1
                                    : undefined
                            }
                        />
                        <SidebarItem
                            to={`/infrastructure/${uuid}/settings`}
                            symbol="settings"
                            name="Settings"
                        />
                        <SidebarItem
                            onClick={() => setShowDeletePopup(true)}
                            symbol="delete"
                            name="Delete"
                            red
                        />
                    </SidebarGroup>
                </Sidebar>
                <div className={styles.side}>
                    <Outlet />
                </div>
                <Popup show={showDeletePopup} onDismiss={dismissDeletePopup}>
                    <Title>
                        Delete{" "}
                        {instance?.display_name ?? instance?.service?.name}?
                    </Title>
                    <Text>
                        Are you sure you want to delete{" "}
                        {instance?.display_name ?? instance?.service?.name}? All
                        data will be permanently deleted.
                    </Text>
                    {deleting && <Progress infinite />}
                    {error && <Error error={error} />}
                    <Horizontal gap={10}>
                        <Spacer />
                        <Button
                            onClick={dismissDeletePopup}
                            disabled={deleting}
                        >
                            Cancel
                        </Button>
                        <Button
                            primary
                            color="red"
                            onClick={onDeleteInstance}
                            disabled={deleting}
                            rightSymbol="delete"
                        >
                            Confirm
                        </Button>
                    </Horizontal>
                </Popup>
            </Horizontal>
        </div>
    );
}
