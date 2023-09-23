import { SubTitle, Title } from "../../components/Text/Text";

import styles from "./SettingsSecurity.module.sass";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import SSHKey, { SSHKeys } from "../../components/SSHKey/SSHKey";
import { APIError } from "../../components/Error/Error";
import ListItem from "../../components/List/ListItem";
import { useFetch } from "../../hooks/useFetch";
import { api } from "../../backend/backend";
import Button from "../../components/Button/Button";
import { ChangeEvent, Fragment, useState } from "react";
import Popup from "../../components/Popup/Popup";
import Spacer from "../../components/Spacer/Spacer";
import Code from "../../components/Code/Code";
import Input from "../../components/Input/Input";
import Progress from "../../components/Progress";

export default function SettingsSecurity() {
    const { data: sshKeys, error } = useFetch<SSHKeys>(api.security.ssh.get);

    const [showPopup, setShowPopup] = useState(false);
    const [authorizedKey, setAuthorizedKey] = useState("");
    const [addError, setAddError] = useState();
    const [adding, setAdding] = useState(false);

    const dismissPopup = () => {
        setShowPopup(false);
        setAuthorizedKey("");
        setAddError(undefined);
    };

    const createSSHKey = () => {
        setAdding(true);
        api.security.ssh
            .add(authorizedKey)
            .then(dismissPopup)
            .catch(setAddError)
            .finally(() => setAdding(false));
    };

    const onAuthorizedKeyChange = (e: ChangeEvent<HTMLInputElement>) => {
        setAuthorizedKey(e.target.value);
    };

    return (
        <Fragment>
            <Vertical gap={20}>
                <Title className={styles.title}>SSH keys</Title>
                <APIError error={error} />
                {!error && sshKeys && (
                    <SSHKeys>
                        {sshKeys?.length === 0 && (
                            <ListItem>No SSH keys found.</ListItem>
                        )}
                        {sshKeys?.map((sshKey) => (
                            <SSHKey
                                key={sshKey.fingerprint_sha_256}
                                type={sshKey.type}
                                fingerprint={sshKey.fingerprint_sha_256}
                            />
                        ))}
                    </SSHKeys>
                )}
                <div>
                    <Button
                        primary
                        leftSymbol="add"
                        onClick={() => setShowPopup(true)}
                    >
                        Create an SSH key
                    </Button>
                </div>
            </Vertical>

            <Popup show={showPopup} onDismiss={dismissPopup}>
                <Title className={styles.title}>Create an SSH key</Title>

                <SubTitle>
                    Step 1: Generate an SSH key if you don't have one
                </SubTitle>
                <Code
                    className={styles.code}
                    code={'ssh-keygen -t ed25519 -C "abc@example.com"'}
                    language={"bash"}
                />

                <SubTitle>Step 2: Paste your public key below</SubTitle>
                <div className={styles.field}>
                    <Input
                        value={authorizedKey}
                        onChange={onAuthorizedKeyChange}
                        placeholder="ssh-ed25519..."
                        disabled={adding}
                    />
                </div>

                <APIError className={styles.error} error={addError} />

                {adding && <Progress infinite />}

                <Horizontal gap={6}>
                    <Spacer />
                    <Button onClick={dismissPopup}>Cancel</Button>
                    <Button primary onClick={createSSHKey}>
                        Create
                    </Button>
                </Horizontal>
            </Popup>
        </Fragment>
    );
}
