import { Vertical } from "../../../components/Layouts/Layouts";
import SSHKey, { SSHKeys } from "../../../components/SSHKey/SSHKey";
import { Errors } from "../../../components/Error/Errors";
import { APIError } from "../../../components/Error/APIError";
import {
    Button,
    ListItem,
    MaterialIcon,
    SelectField,
    SelectOption,
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
        error: keysError,
        isLoading: isKeysLoading,
    } = useQuery({
        queryKey: ["ssh_keys"],
        queryFn: api.security.ssh.get,
    });

    const {
        data: sshUsers,
        error: usersError,
        isLoading: isUsersLoading,
    } = useQuery({
        queryKey: ["ssh_users"],
        queryFn: api.security.ssh.users,
    });

    const [showPopup, setShowPopup] = useState(false);
    const [authorizedKey, setAuthorizedKey] = useState("");
    const [username, setUsername] = useState("");
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
            .add(authorizedKey, username)
            .then(() => {
                dismissPopup();
                queryClient.invalidateQueries({
                    queryKey: ["ssh_keys"],
                });
            })
            .catch(setAddError)
            .finally(() => setAdding(false));
    };

    const deleteSSHKey = (fingerprint: string, username: string) => {
        api.security.ssh
            .delete(fingerprint, username)
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

    const onUsernameChange = (value: any) => {
        setUsername(value);
    };

    const popupActions = (
        <Fragment>
            <Button onClick={dismissPopup}>Cancel</Button>
            <Button
                variant="colored"
                disabled={authorizedKey === "" || username === ""}
                onClick={createSSHKey}
            >
                Create
            </Button>
        </Fragment>
    );

    const error = keysError || usersError;
    const isLoading = isKeysLoading || isUsersLoading;

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
                                username={sshKey.username}
                                onDelete={() =>
                                    deleteSSHKey(
                                        sshKey.fingerprint_sha_256,
                                        sshKey.username
                                    )
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
                title="Add SSH key"
                actions={popupActions}
            >
                <TextField
                    id="authorized-key"
                    value={authorizedKey}
                    label="Public key"
                    description="Paste your public SSH key here. It should start with ssh-ed25519."
                    required
                    onChange={onAuthorizedKeyChange}
                    placeholder="ssh-ed25519..."
                    disabled={adding}
                />
                <SelectField
                    onChange={onUsernameChange}
                    value={username}
                    label="User"
                    description="Select the user to associate this key with."
                    required
                >
                    {sshUsers?.map((user) => (
                        <SelectOption key={user} value={user}>
                            {user}
                        </SelectOption>
                    ))}
                </SelectField>

                <APIError error={addError} />
                {adding && <Progress infinite />}
            </Popup>
        </Content>
    );
}
