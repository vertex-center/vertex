import { useState } from "react";
import { useParams } from "react-router-dom";
import { Horizontal, Vertical } from "../../../../components/Layouts/Layouts";
import { Button, MaterialIcon, Title } from "@vertex-center/components";
import { api } from "../../../../backend/api/backend";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../../components/KeyValueInfo/KeyValueInfo";
import byteSize from "byte-size";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useDockerInfo } from "../../hooks/useContainer";
import Content from "../../../../components/Content/Content";

export default function ContainerDocker() {
    const { uuid } = useParams();

    const { dockerInfo: info, isLoading, error } = useDockerInfo(uuid);

    const [recreatingContainer, setRecreatingContainer] = useState(false);
    const [recreatingContainerError, setRecreatingContainerError] = useState();

    const recreateContainer = () => {
        setRecreatingContainer(true);
        setRecreatingContainerError(undefined);

        api.vxContainers
            .container(uuid)
            .docker.recreate()
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

    if (error) return <APIError error={error} />;

    return (
        <Content>
            <ProgressOverlay show={isLoading} />
            <APIError error={recreatingContainerError} />

            <Title variant="h3">Container</Title>
            <KeyValueGroup>
                <KeyValueInfo
                    name="ID"
                    type="code"
                    icon="tag"
                    loading={isLoading}
                >
                    {info?.container?.id}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Container Name"
                    type="code"
                    icon="badge"
                    loading={isLoading}
                >
                    {info?.container?.name}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Platform"
                    type="code"
                    icon="computer"
                    loading={isLoading}
                >
                    {info?.container?.platform}
                </KeyValueInfo>
            </KeyValueGroup>

            <Title variant="h3">Image</Title>
            <KeyValueGroup>
                <KeyValueInfo
                    name="ID"
                    type="code"
                    icon="tag"
                    loading={isLoading}
                >
                    {info?.image?.id}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Architecture"
                    type="code"
                    icon="memory"
                    loading={isLoading}
                >
                    {info?.image?.architecture}
                </KeyValueInfo>
                <KeyValueInfo
                    name="OS"
                    type="code"
                    icon="computer"
                    loading={isLoading}
                >
                    {info?.image?.os}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Size"
                    type="code"
                    icon="straighten"
                    loading={isLoading}
                >
                    {byteSize(info?.image?.size).toString()}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Tags"
                    type="code"
                    icon="sell"
                    loading={isLoading}
                >
                    {info?.image?.tags?.join(", ")}
                </KeyValueInfo>
            </KeyValueGroup>

            <Title variant="h3">Actions</Title>
            <Vertical gap={8}>
                <Horizontal alignItems="center" gap={20}>
                    <Button
                        leftIcon={<MaterialIcon icon="restart_alt" />}
                        onClick={recreateContainer}
                        disabled={recreatingContainer ?? isLoading}
                    >
                        Recreate container
                    </Button>
                </Horizontal>
            </Vertical>
        </Content>
    );
}
