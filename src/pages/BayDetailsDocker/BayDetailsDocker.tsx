import { useState } from "react";
import { Text, Title } from "../../components/Text/Text";

import styles from "./BayDetailsDocker.module.sass";
import { useParams } from "react-router-dom";
import Loading from "../../components/Loading/Loading";
import useInstance from "../../hooks/useInstance";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import Spacer from "../../components/Spacer/Spacer";
import Button from "../../components/Button/Button";
import { Error } from "../../components/Error/Error";
import { useFetch } from "../../hooks/useFetch";
import { DockerContainerInfo } from "../../models/docker";
import { api } from "../../backend/backend";

export default function BayDetailsDocker() {
    const { uuid } = useParams();

    const { instance } = useInstance(uuid);

    const { data: containerInfo, reload: reloadContainerInfo } =
        useFetch<DockerContainerInfo>(() => api.instance.docker.get(uuid));

    const [recreatingContainer, setRecreatingContainer] = useState(false);
    const [recreatingContainerError, setRecreatingContainerError] = useState();

    let label;
    if (instance?.install_method === "docker") {
        label = <span className={styles.enabled}>enabled</span>;
    } else {
        label = <span className={styles.disabled}>disabled</span>;
    }

    const recreateContainer = () => {
        setRecreatingContainer(true);
        setRecreatingContainerError(undefined);

        api.instance.docker
            .recreate(uuid)
            .then(() => {
                reloadContainerInfo().then();
            })
            .catch((err) => {
                setRecreatingContainerError(
                    err?.response?.data?.message ?? err?.message
                );
                console.error(err);
            })
            .finally(() => {
                setRecreatingContainer(false);
            });
    };

    return (
        <Vertical gap={40}>
            <Vertical gap={20}>
                <Title className={styles.title}>Docker</Title>
                {!instance && <Loading />}
                {instance && (
                    <Text className={styles.content}>
                        Docker is {label} for this instance.
                    </Text>
                )}
            </Vertical>

            <Vertical gap={20}>
                <Title className={styles.title}>Details</Title>
                <Vertical className={styles.content} gap={8}>
                    <Horizontal gap={12} alignItems="center">
                        <Text>Container ID</Text>
                        <Spacer />
                        <code>{containerInfo?.id}</code>
                    </Horizontal>
                    <Horizontal gap={12} alignItems="center">
                        <Text>Container Name</Text>
                        <Spacer />
                        <code>{containerInfo?.name}</code>
                    </Horizontal>
                    <Horizontal gap={12} alignItems="center">
                        <Text>Image</Text>
                        <Spacer />
                        <code>{containerInfo?.image}</code>
                    </Horizontal>
                    <Horizontal gap={12} alignItems="center">
                        <Text>Platform</Text>
                        <Spacer />
                        <code>{containerInfo?.platform}</code>
                    </Horizontal>
                </Vertical>
            </Vertical>

            <Vertical gap={20}>
                <Title className={styles.title}>Actions</Title>
                <Vertical className={styles.content} gap={8}>
                    <Horizontal alignItems="center" gap={20}>
                        <Button
                            leftSymbol="restart_alt"
                            onClick={recreateContainer}
                            disabled={recreatingContainer}
                        >
                            Recreate container
                        </Button>
                        <Error error={recreatingContainerError} />
                    </Horizontal>
                </Vertical>
            </Vertical>
        </Vertical>
    );
}
