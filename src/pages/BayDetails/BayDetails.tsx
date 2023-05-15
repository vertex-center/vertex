import Bay from "../../components/Bay/Bay";
import { useEffect, useState } from "react";
import {
    deleteInstance,
    route,
    startInstance,
    stopInstance,
} from "../../backend/backend";
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
import Header from "../../components/Header/Header";
import Sidebar, { SidebarItem } from "../../components/Sidebar/Sidebar";
import Popup from "../../components/Popup/Popup";
import { Text, Title } from "../../components/Text/Text";
import Button from "../../components/Button/Button";
import Progress from "../../components/Progress";
import { SiDocker } from "@icons-pack/react-simple-icons";
import useInstance from "../../hooks/useInstance";

export default function BayDetails() {
    const { uuid } = useParams();
    const navigate = useNavigate();

    const { instance, setInstance } = useInstance(uuid);

    const [showDeletePopup, setShowDeletePopup] = useState<boolean>();
    const [deleting, setDeleting] = useState<boolean>(false);

    useEffect(() => {
        if (uuid === undefined) return;

        const sse = registerSSE(route(`/instance/${uuid}/events`));

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
            await startInstance(uuid);
        } else {
            await stopInstance(uuid);
        }
    };

    const onDeleteInstance = () => {
        setDeleting(true);
        deleteInstance(uuid).then(() => {
            navigate("/");
        });
    };

    const dismissDeletePopup = () => {
        setShowDeletePopup(false);
    };

    return (
        <div className={styles.details}>
            <Header />
            <div className={styles.bay}>
                <Bay
                    showCables
                    instances={[
                        {
                            name: instance?.name,
                            status: instance?.status,
                            onPower: () => toggleInstance(uuid),
                            use_docker: instance?.use_docker,
                        },
                    ]}
                />
            </div>
            <Horizontal className={styles.content}>
                <Sidebar>
                    <SidebarItem to="/" symbol="arrow_back" name="Back" />
                    <SidebarItem
                        to={`/bay/${uuid}/`}
                        symbol="home"
                        name="Home"
                    />
                    <div className={styles.separator} />
                    <SidebarItem
                        to={`/bay/${uuid}/logs`}
                        symbol="terminal"
                        name="Logs"
                    />
                    {/* Uptime status is disabled for now */}
                    {/*<SidebarItem*/}
                    {/*    to={`/bay/${uuid}/status`}*/}
                    {/*    symbol="vital_signs"*/}
                    {/*    name="Status"*/}
                    {/*/>*/}
                    <SidebarItem
                        to={`/bay/${uuid}/environment`}
                        symbol="tune"
                        name="Environment"
                    />
                    {instance?.use_docker && (
                        <SidebarItem
                            to={`/bay/${uuid}/docker`}
                            symbol={<SiDocker size={20} />}
                            name="Docker"
                        />
                    )}
                    {!instance?.use_docker && (
                        <SidebarItem
                            to={`/bay/${uuid}/dependencies`}
                            symbol="widgets"
                            name="Dependencies"
                        />
                    )}
                    <SidebarItem
                        to={`/bay/${uuid}/settings`}
                        symbol="settings"
                        name="Settings"
                    />
                    <Spacer />
                    <SidebarItem
                        onClick={() => setShowDeletePopup(true)}
                        symbol="delete"
                        name="Delete"
                        red
                    />
                </Sidebar>
                <div className={styles.side}>
                    <Outlet />
                </div>
                <Popup show={showDeletePopup} onDismiss={dismissDeletePopup}>
                    <Title>Delete {instance?.name}?</Title>
                    <Text>
                        Are you sure you want to delete {instance?.name}? All
                        data will be permanently deleted.
                    </Text>
                    {deleting && <Progress infinite />}
                    <Horizontal gap={12}>
                        <Spacer />
                        <Button onClick={onDeleteInstance} disabled={deleting}>
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
