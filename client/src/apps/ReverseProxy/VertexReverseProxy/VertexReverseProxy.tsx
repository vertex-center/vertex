import { api } from "../../../backend/api/backend";
import React, { useState } from "react";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import styles from "./VertexReverseProxy.module.sass";
import { APIError } from "../../../components/Error/APIError";
import ProxyRedirect from "../../../components/ProxyRedirect/ProxyRedirect";
import { Horizontal, Vertical } from "../../../components/Layouts/Layouts";
import {
    Button,
    FormItem,
    Input,
    MaterialIcon,
    Title,
} from "@vertex-center/components";
import Popup, { PopupActions } from "../../../components/Popup/Popup";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import NoItems from "../../../components/NoItems/NoItems";
import Content from "../../../components/Content/Content";
import { Plus, ShareNetwork } from "@phosphor-icons/react";

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
        onSettled: async () => {
            await queryClient.invalidateQueries({
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
        onSettled: async () => {
            await queryClient.invalidateQueries({
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
        <Content>
            <ProgressOverlay show={isLoading} />
            <Title variant="h2">Vertex Reverse Proxy</Title>

            <Vertical gap={20}>
                <div className={styles.redirects}>
                    {error && <APIError error={error} />}
                    {!error &&
                        Object.entries(redirects ?? {}).map(
                            ([uuid, redirect]) => (
                                <ProxyRedirect
                                    key={uuid}
                                    enabled={true}
                                    source={redirect.source}
                                    target={redirect.target}
                                    onDelete={() => mutationDelete.mutate(uuid)}
                                />
                            )
                        )}
                    {!error &&
                        redirects &&
                        Object.keys(redirects).length === 0 && (
                            <NoItems
                                text="No redirections found."
                                icon={<ShareNetwork />}
                            />
                        )}
                </div>
                <Horizontal gap={10}>
                    <Button
                        variant="colored"
                        onClick={openNewRedirectPopup}
                        leftIcon={<Plus />}
                    >
                        Add redirection
                    </Button>
                </Horizontal>
            </Vertical>
            {showNewRedirectPopup && (
                <Popup
                    onDismiss={closeNewRedirectPopup}
                    title="New redirection"
                >
                    <FormItem label="Source">
                        <Input
                            className={styles.input}
                            value={source}
                            onChange={onSourceChange}
                        />
                    </FormItem>
                    <FormItem label="Target">
                        <Input value={target} onChange={onTargetChange} />
                    </FormItem>
                    <PopupActions>
                        <Button onClick={closeNewRedirectPopup}>Cancel</Button>
                        <Button
                            variant="colored"
                            onClick={async () =>
                                mutationAdd.mutate({ source, target })
                            }
                            rightIcon={<MaterialIcon icon="send" />}
                        >
                            Send
                        </Button>
                    </PopupActions>
                </Popup>
            )}
        </Content>
    );
}
