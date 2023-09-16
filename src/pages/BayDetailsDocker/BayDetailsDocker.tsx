import { useState } from "react";
import { Title } from "../../components/Text/Text";

import styles from "./BayDetailsDocker.module.sass";
import { useParams } from "react-router-dom";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import Button from "../../components/Button/Button";
import { Error } from "../../components/Error/Error";
import { useFetch } from "../../hooks/useFetch";
import { DockerContainerInfo } from "../../models/docker";
import { api } from "../../backend/backend";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../components/KeyValueInfo/KeyValueInfo";
import byteSize from "byte-size";

export default function BayDetailsDocker() {
    const { uuid } = useParams();

    const { data: info, reload: reloadContainerInfo } =
        useFetch<DockerContainerInfo>(() => api.instance.docker.get(uuid));

    const [recreatingContainer, setRecreatingContainer] = useState(false);
    const [recreatingContainerError, setRecreatingContainerError] = useState();

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
        <Vertical gap={30}>
            <Vertical gap={20}>
                <Title className={styles.title}>Container</Title>
                <KeyValueGroup>
                    <KeyValueInfo name="ID" type="code" symbol="tag">
                        {info?.container?.id}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Container Name"
                        type="code"
                        symbol="badge"
                    >
                        {info?.container?.name}
                    </KeyValueInfo>
                    <KeyValueInfo name="Platform" type="code" symbol="computer">
                        {info?.container?.platform}
                    </KeyValueInfo>
                </KeyValueGroup>
            </Vertical>

            <Vertical gap={20}>
                <Title className={styles.title}>Image</Title>
                <KeyValueGroup>
                    <KeyValueInfo name="ID" type="code" symbol="tag">
                        {info?.image?.id}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Architecture"
                        type="code"
                        symbol="memory"
                    >
                        {info?.image?.architecture}
                    </KeyValueInfo>
                    <KeyValueInfo name="OS" type="code" symbol="computer">
                        {info?.image?.os}
                    </KeyValueInfo>
                    <KeyValueInfo name="Size" type="code" symbol="straighten">
                        {byteSize(info?.image?.size).toString()}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Virtual size"
                        type="code"
                        symbol="straighten"
                    >
                        {byteSize(info?.image?.virtual_size).toString()}
                    </KeyValueInfo>
                    <KeyValueInfo name="Tags" type="code" symbol="sell">
                        {info?.image?.tags?.join(", ")}
                    </KeyValueInfo>
                </KeyValueGroup>
            </Vertical>

            <Vertical gap={20}>
                <Title className={styles.title}>Actions</Title>
                <Vertical gap={8}>
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
