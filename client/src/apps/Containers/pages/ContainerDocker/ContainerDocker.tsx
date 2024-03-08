import { useState } from "react";
import { useParams } from "react-router-dom";
import { Button, Horizontal, Title, Vertical } from "@vertex-center/components";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../../components/KeyValueInfo/KeyValueInfo";
import byteSize from "byte-size";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useDockerInfo } from "../../hooks/useContainer";
import Content from "../../../../components/Content/Content";
import { API } from "../../backend/api";
import {
    AppWindow,
    ArrowClockwise,
    ArrowsOut,
    Cpu,
    IdentificationBadge,
    Tag,
} from "@phosphor-icons/react";
import { useReloadContainer } from "../../hooks/useContainers";

export default function ContainerDocker() {
    const { uuid } = useParams();

    const {
        dockerInfo: info,
        isLoading,
        error: errorDockerInfo,
    } = useDockerInfo(uuid);

    const [recreatingContainer, setRecreatingContainer] = useState(false);
    const [recreatingContainerError, setRecreatingContainerError] = useState();

    const recreateContainer = () => {
        setRecreatingContainer(true);
        setRecreatingContainerError(undefined);

        API.recreateDocker(uuid)
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

    const {
        reloadContainer,
        isPending,
        error: reloadContainerError,
    } = useReloadContainer();

    const error =
        errorDockerInfo || recreatingContainerError || reloadContainerError;

    if (error) return <APIError error={error} />;

    return (
        <Content>
            <Title variant="h3">Container</Title>
            <KeyValueGroup>
                <KeyValueInfo
                    name="ID"
                    type="code"
                    icon={<Tag />}
                    loading={isLoading}
                >
                    {info?.container?.id}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Container Name"
                    type="code"
                    icon={<IdentificationBadge />}
                    loading={isLoading}
                >
                    {info?.container?.name}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Platform"
                    type="code"
                    icon={<Cpu />}
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
                    icon={<Tag />}
                    loading={isLoading}
                >
                    {info?.image?.id}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Architecture"
                    type="code"
                    icon={<Cpu />}
                    loading={isLoading}
                >
                    {info?.image?.architecture}
                </KeyValueInfo>
                <KeyValueInfo
                    name="OS"
                    type="code"
                    icon={<AppWindow />}
                    loading={isLoading}
                >
                    {info?.image?.os}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Size"
                    type="code"
                    icon={<ArrowsOut />}
                    loading={isLoading}
                >
                    {byteSize(info?.image?.size).toString()}
                </KeyValueInfo>
                <KeyValueInfo
                    name="Tags"
                    type="code"
                    icon={<Tag />}
                    loading={isLoading}
                >
                    {info?.image?.tags?.join(", ")}
                </KeyValueInfo>
            </KeyValueGroup>

            <Title variant="h3">Actions</Title>
            <Vertical gap={8}>
                <Horizontal alignItems="center" gap={20}>
                    <Button
                        leftIcon={<ArrowClockwise />}
                        onClick={recreateContainer}
                        disabled={recreatingContainer || isLoading}
                    >
                        Recreate container
                    </Button>
                    <Button
                        leftIcon={<ArrowClockwise />}
                        onClick={() => reloadContainer(uuid)}
                        disabled={isPending || isLoading}
                    >
                        Reload status
                    </Button>
                </Horizontal>
            </Vertical>

            <ProgressOverlay show={isLoading || isPending} />
        </Content>
    );
}
