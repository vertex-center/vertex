import { Vertical } from "../../../components/Layouts/Layouts";
import SSHKey, { SSHKeys } from "../../../components/SSHKey/SSHKey";
import { Errors } from "../../../components/Error/Errors";
import { APIError } from "../../../components/Error/APIError";
import {
    Button,
    Code,
    ListItem,
    MaterialIcon,
    TextField,
    Title,
} from "@vertex-center/components";
import { api } from "../../../backend/api/backend";
import { ChangeEvent, Fragment, useState } from "react";
import Popup from "../../../components/Popup/Popup";
import Progress from "../../../components/Progress";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import Content from "../../../components/Content/Content";

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

    const popupActions = (
        <Fragment>
            <Button onClick={dismissPopup}>Cancel</Button>
            <Button
                variant="colored"
                disabled={authorizedKey === ""}
                onClick={createSSHKey}
            >
                Create
            </Button>
        </Fragment>
    );

    return (
        <Content>
            <ProgressOverlay show={isLoading} />
            <Title variant="h2">SSH keys</Title>
            <Vertical gap={20}>
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

            <Popup
                show={showPopup}
                onDismiss={dismissPopup}
                title="Create SSH key"
                actions={popupActions}
            >
                <Title variant="h4">
                    Step 1: Generate an SSH key if you don't have one
                </Title>
                <Code language={"bash"}>
                    ssh-keygen -t ed25519 -C "abc@example.com"
                </Code>

                <Title variant="h4">Step 2: Paste your public key below</Title>
                <TextField
                    id="authorized-key"
                    value={authorizedKey}
                    onChange={onAuthorizedKeyChange}
                    placeholder="ssh-ed25519..."
                    disabled={adding}
                />

                <APIError error={addError} />
                {adding && <Progress infinite />}
            </Popup>
        </Content>
    );
}
