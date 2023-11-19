import Container, { Containers } from "../Container/Container";
import { APIError } from "../Error/APIError";
import { Fragment, useEffect, useState } from "react";
import { ProgressOverlay } from "../Progress/Progress";
import { Container as ContainerModel } from "../../models/container";
import { api } from "../../backend/api/backend";
import { useServerEvent } from "../../hooks/useEvent";
import { useQueryClient } from "@tanstack/react-query";
import { useContainers } from "../../apps/Containers/hooks/useContainers";

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
        if (!containers) return;
        const inst = Object.values(containers ?? {})?.[0];
        if (!inst) {
            setContainer({
                display_name: name,
                status: "not-installed",
            });
            return;
        }
        setContainer({
            ...inst,
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

    const onPower = async (inst: Inst) => {
        if (!inst) {
            console.error("Container not found");
            return;
        }
        if (inst?.status === "off" || inst?.status === "error") {
            await api.vxContainers.container(inst.uuid).start();
            return;
        }
        await api.vxContainers.container(inst.uuid).stop();
    };

    const route = container?.uuid
        ? `/app/vx-containers/container/${container?.uuid}/events`
        : "";

    useServerEvent(route, {
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
                        to: container?.uuid
                            ? `/app/vx-containers/${container?.uuid}`
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
