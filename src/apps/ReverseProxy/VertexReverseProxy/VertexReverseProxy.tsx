import { api } from "../../../backend/api/backend";
import React, { useState } from "react";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import styles from "./VertexReverseProxy.module.sass";
import { Title } from "../../../components/Text/Text";
import { APIError } from "../../../components/Error/APIError";
import ProxyRedirect from "../../../components/ProxyRedirect/ProxyRedirect";
import { Horizontal, Vertical } from "../../../components/Layouts/Layouts";
import Button from "../../../components/Button/Button";
import Popup from "../../../components/Popup/Popup";
import Input from "../../../components/Input/Input";
import Spacer from "../../../components/Spacer/Spacer";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export default function VertexReverseProxy() {
    const queryClient = useQueryClient();
    const {
        data: redirects,
        error,
        isLoading,
    } = useQuery({
        queryKey: ["redirects"],
        queryFn: api.vxReverseProxy.redirects.all,
    });

    const mutationDelete = useMutation({
        mutationFn: api.vxReverseProxy.redirects.delete,
        onSuccess: () => {
            closeNewRedirectPopup();
        },
        onSettled: () => {
            queryClient.invalidateQueries({
                queryKey: ["redirects"],
            });
        },
    });

    const mutationAdd = useMutation({
        mutationFn: ({ source, target }: { source: string; target: string }) =>
            api.vxReverseProxy.redirects.add(source, target),
        onSuccess: () => {
            closeNewRedirectPopup();
        },
        onSettled: () => {
            queryClient.invalidateQueries({
                queryKey: ["redirects"],
            });
        },
    });

    const [showNewRedirectPopup, setShowNewRedirectPopup] = useState(false);

    const [source, setSource] = useState("");
    const [target, setTarget] = useState("");

    const onSourceChange = (e: any) => setSource(e.target.value);
    const onTargetChange = (e: any) => setTarget(e.target.value);

    const openNewRedirectPopup = () => setShowNewRedirectPopup(true);
    const closeNewRedirectPopup = () => setShowNewRedirectPopup(false);

    return (
        <Vertical gap={20}>
            <ProgressOverlay show={isLoading} />
            <Title className={styles.title}>Vertex Reverse Proxy</Title>

            <div className={styles.redirects}>
                {error && <APIError error={error} />}
                {!error &&
                    Object.entries(redirects ?? {}).map(([uuid, redirect]) => (
                        <ProxyRedirect
                            key={uuid}
                            enabled={true}
                            source={redirect.source}
                            target={redirect.target}
                            onDelete={() => mutationDelete.mutate(uuid)}
                        />
                    ))}
            </div>
            <Horizontal className={styles.addRedirect} gap={10}>
                <Button primary onClick={openNewRedirectPopup} leftIcon="add">
                    Add redirection
                </Button>
            </Horizontal>
            <Popup
                show={showNewRedirectPopup}
                onDismiss={closeNewRedirectPopup}
            >
                <Title>New redirection</Title>
                <Vertical gap={20} className={styles.input}>
                    <Input
                        className={styles.input}
                        label="Source"
                        value={source}
                        onChange={onSourceChange}
                    />
                    <Input
                        label="Target"
                        value={target}
                        onChange={onTargetChange}
                    />
                </Vertical>
                <Horizontal gap={10}>
                    <Spacer />
                    <Button onClick={closeNewRedirectPopup}>Cancel</Button>
                    <Button
                        primary
                        onClick={async () =>
                            mutationAdd.mutate({ source, target })
                        }
                    >
                        Send
                    </Button>
                </Horizontal>
            </Popup>
        </Vertical>
    );
}
