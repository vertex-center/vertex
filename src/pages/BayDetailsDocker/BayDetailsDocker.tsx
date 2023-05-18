import { useEffect, useState } from "react";
import { Text, Title } from "../../components/Text/Text";

import styles from "./BayDetailsDocker.module.sass";
import { useParams } from "react-router-dom";
import Loading from "../../components/Loading/Loading";
import useInstance from "../../hooks/useInstance";
import {
    DockerContainerInfo,
    getInstanceDockerContainerInfo,
} from "../../backend/backend";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import Spacer from "../../components/Spacer/Spacer";

export default function BayDetailsDocker() {
    const { uuid } = useParams();

    const { instance } = useInstance(uuid);

    const [containerInfo, setContainerInfo] = useState<DockerContainerInfo>();

    let label;

    useEffect(() => {
        getInstanceDockerContainerInfo(uuid)
            .then((info) => {
                console.log(info);
                setContainerInfo(info);
            })
            .catch();
    }, [uuid]);

    if (instance?.use_docker) {
        label = <span className={styles.enabled}>enabled</span>;
    } else {
        label = <span className={styles.disabled}>disabled</span>;
    }

    return (
        <Vertical gap={40}>
            <Vertical gap={20}>
                <Title>Docker</Title>
                {!instance && <Loading />}
                {instance && <Text>Docker is {label} for this instance.</Text>}
            </Vertical>

            <Vertical gap={20}>
                <Title>Details</Title>
                <Vertical gap={8}>
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
        </Vertical>
    );
}