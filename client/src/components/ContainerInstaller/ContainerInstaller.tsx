import Container, { Containers } from "../Container/Container";
import { APIError } from "../Error/APIError";
import { Fragment, useEffect, useState } from "react";
import { ProgressOverlay } from "../Progress/Progress";
import { Container as ContainerModel } from "../../apps/Containers/backend/models";
import { useServerEvent } from "../../hooks/useEvent";
import { useQueryClient } from "@tanstack/react-query";
import { useContainers } from "../../apps/Containers/hooks/useContainers";
import { API } from "../../apps/Containers/backend/api";

type Props = {
    name: string;
    tag: string;
    install: () => Promise<any>;
};

type Inst = Partial<
    ContainerModel & {
        onPower: (inst: ContainerModel) => Promise<void>;
    }
>;

export default function ContainerInstaller(props: Readonly<Props>) {
    const { name, tag, install } = props;
    const queryClient = useQueryClient();

    const {
        containers,
        isLoading: isLoadingContainers,
        error: errorContainers,
    } = useContainers({
        tags: [tag],
    });

    const [downloading, setDownloading] = useState(false);
    const [error, setError] = useState();

    const [container, setContainer] = useState<Inst>();

    useEffect(() => {
        if (!containers || containers.length === 0) {
            setContainer({
                name: name,
                status: "not-installed",
            });
            return;
        }
        setContainer({
            ...containers[0],
            onPower: onPower,
        });
    }, [containers]);

    const onInstall = () => {
        setError(undefined);
        setDownloading(true);

        install()
            .catch(setError)
            .finally(() => {
                setDownloading(false);
                queryClient.invalidateQueries({
                    queryKey: ["containers"],
                });
            });
    };

    const onPower = async (c: Inst) => {
        if (!c) {
            console.error("Container not found");
            return;
        }
        if (c?.status === "off" || c?.status === "error") {
            await API.startContainer(c.id);
            return;
        }
        await API.stopContainer(c.id);
    };

    const route = container?.id ? `/container/${container?.id}/events` : "";

    // @ts-ignore
    useServerEvent(window.api_urls.containers, route, {
        status_change: (e) => {
            setContainer((container) => ({ ...container, status: e.data }));
        },
    });

    return (
        <Fragment>
            <Containers>
                <Container
                    container={{
                        value: container,
                        to: container?.id
                            ? `/containers/${container?.id}/logs`
                            : undefined,
                        onInstall: onInstall,
                        onPower: () => onPower(container),
                    }}
                />
            </Containers>
            <ProgressOverlay show={downloading || isLoadingContainers} />
            <APIError error={error ?? errorContainers} />
        </Fragment>
    );
}
