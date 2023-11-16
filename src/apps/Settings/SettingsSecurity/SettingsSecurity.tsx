import { SubTitle, Title } from "../../../components/Text/Text";

import styles from "./SettingsSecurity.module.sass";
import { Horizontal, Vertical } from "../../../components/Layouts/Layouts";
import SSHKey, { SSHKeys } from "../../../components/SSHKey/SSHKey";
import { Errors } from "../../../components/Error/Errors";
import { APIError } from "../../../components/Error/APIError";
import {
    Button,
    ListItem,
    MaterialIcon,
    TextField,
} from "@vertex-center/components";
import { api } from "../../../backend/api/backend";
import { ChangeEvent, Fragment, useState } from "react";
import Popup from "../../../components/Popup/Popup";
import Spacer from "../../../components/Spacer/Spacer";
import Code from "../../../components/Code/Code";
import Progress from "../../../components/Progress";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useQuery, useQueryClient } from "@tanstack/react-query";

export default function SettingsSecurity() {
    const queryClient = useQueryClient();
    const {
        data: sshKeys,
        error,
        isLoading,
    } = useQuery({
        queryKey: ["ssh_keys"],
        queryFn: api.security.ssh.get,
    });

    const [showPopup, setShowPopup] = useState(false);
    const [authorizedKey, setAuthorizedKey] = useState("");
    const [addError, setAddError] = useState();
    const [deleteError, setDeleteError] = useState();
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
            .then(() => {
                dismissPopup();
                queryClient.invalidateQueries({
                    queryKey: ["ssh_keys"],
                });
            })
            .catch(setAddError)
            .finally(() => setAdding(false));
    };

    const deleteSSHKey = (fingerprint: string) => {
        api.security.ssh
            .delete(fingerprint)
            .then(() => {
                queryClient.invalidateQueries({
                    queryKey: ["ssh_keys"],
                });
            })
            .catch(setDeleteError);
    };

    const onAuthorizedKeyChange = (e: ChangeEvent<HTMLInputElement>) => {
        setAuthorizedKey(e.target.value);
    };

    return (
        <Fragment>
            <ProgressOverlay show={isLoading} />
            <Vertical gap={20}>
                <Title className={styles.title}>SSH keys</Title>
                {(error || deleteError) && (
                    <Errors>
                        <APIError error={error} />
                        <APIError error={deleteError} />
                    </Errors>
                )}
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
                                onDelete={() =>
                                    deleteSSHKey(sshKey.fingerprint_sha_256)
                                }
                            />
                        ))}
                    </SSHKeys>
                )}
                <div>
                    <Button
                        variant="colored"
                        leftIcon={<MaterialIcon icon="add" />}
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
                    <TextField
                        id="authorized-key"
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
                    <Button
                        variant="colored"
                        disabled={authorizedKey === ""}
                        onClick={createSSHKey}
                    >
                        Create
                    </Button>
                </Horizontal>
            </Popup>
        </Fragment>
    );
}
