import { BigTitle, Title } from "../../../components/Text/Text";
import styles from "./ReverseProxyApp.module.sass";
import ProxyRedirect from "../../../components/ProxyRedirect/ProxyRedirect";
import { Horizontal, Vertical } from "../../../components/Layouts/Layouts";
import Button from "../../../components/Button/Button";
import { useFetch } from "../../../hooks/useFetch";
import { api } from "../../../backend/backend";
import Popup from "../../../components/Popup/Popup";
import React, { useState } from "react";
import Input from "../../../components/Input/Input";
import Spacer from "../../../components/Spacer/Spacer";
import { APIError } from "../../../components/Error/Error";

export default function ReverseProxyApp() {
    const {
        data: redirects,
        error,
        reload,
    } = useFetch<ProxyRedirects>(api.proxy.redirects.get);

    const [showNewRedirectPopup, setShowNewRedirectPopup] = useState(false);

    const [source, setSource] = useState("");
    const [target, setTarget] = useState("");

    const onSourceChange = (e: any) => setSource(e.target.value);
    const onTargetChange = (e: any) => setTarget(e.target.value);

    const openNewRedirectPopup = () => setShowNewRedirectPopup(true);
    const closeNewRedirectPopup = () => setShowNewRedirectPopup(false);

    const addNewRedirection = () => {
        api.proxy.redirects
            .add(source, target)
            .then(reload)
            .catch(console.error)
            .finally(closeNewRedirectPopup);
    };

    const onDelete = (uuid: string) => {
        api.proxy.redirects.delete(uuid).then(reload).catch(console.error);
    };

    return (
        <div>
            <div className={styles.title}>
                <BigTitle>Reverse Proxy</BigTitle>
            </div>

            <div className={styles.redirects}>
                {error && <APIError error={error} />}
                {!error &&
                    Object.entries(redirects ?? {}).map(([uuid, redirect]) => (
                        <ProxyRedirect
                            enabled={true}
                            source={redirect.source}
                            target={redirect.target}
                            onDelete={() => onDelete(uuid)}
                        />
                    ))}
            </div>
            <Horizontal className={styles.addRedirect} gap={10}>
                <Button primary onClick={openNewRedirectPopup} leftSymbol="add">
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
                    <Button primary onClick={addNewRedirection}>
                        Send
                    </Button>
                </Horizontal>
            </Popup>
        </div>
    );
}
