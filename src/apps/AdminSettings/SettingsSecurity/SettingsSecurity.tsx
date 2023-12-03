import { Vertical } from "../../../components/Layouts/Layouts";
import SSHKey from "../../../components/SSHKey/SSHKey";
import { APIError } from "../../../components/Error/APIError";
import {
    Button,
    List,
    ListItem,
    MaterialIcon,
    SelectField,
    SelectOption,
    TextField,
    Title,
} from "@vertex-center/components";
import { ChangeEvent, Fragment, useState } from "react";
import Popup from "../../../components/Popup/Popup";
import Progress from "../../../components/Progress";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useQueryClient } from "@tanstack/react-query";
import Content from "../../../components/Content/Content";
import { useSSHKeys } from "../hooks/useSSHKeys";
import { useSSHUsers } from "../hooks/useSSHUsers";
import { useCreateSSHKey } from "../hooks/useCreateSSHKey";
import { useDeleteSSHKey } from "../hooks/useDeleteSSHKey";

export default function SettingsSecurity() {
    const queryClient = useQueryClient();

    const { sshKeys, keysError, isKeysLoading } = useSSHKeys();
    const { sshUsers, sshUsersError, isSSHUsersLoading } = useSSHUsers();

    const [showPopup, setShowPopup] = useState(false);
    const [authorizedKey, setAuthorizedKey] = useState("");
    const [username, setUsername] = useState("");

    const dismissPopup = () => {
        setShowPopup(false);
        setAuthorizedKey("");
        resetCreateKey();
    };

    const { createKey, isCreatingKey, errorCreateKey, resetCreateKey } =
        useCreateSSHKey({
            onSuccess: () => {
                dismissPopup();
                queryClient.invalidateQueries({
                    queryKey: ["admin_ssh_keys"],
                });
            },
        });

    const { deleteKey, isDeletingKey, errorDeleteKey } = useDeleteSSHKey({
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: ["admin_ssh_keys"],
            });
        },
    });

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
                onClick={() =>
                    createKey({ authorized_key: authorizedKey, username })
                }
            >
                Create
            </Button>
        </Fragment>
    );

    const error = keysError || sshUsersError || errorDeleteKey;
    const isLoading = isKeysLoading || isSSHUsersLoading || isDeletingKey;

    return (
        <Content>
            <ProgressOverlay show={isLoading} />
            <Title variant="h2">SSH keys</Title>
            <Vertical gap={20}>
                <APIError error={error} />
                {!error && sshKeys && (
                    <List>
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
                                    deleteKey({
                                        fingerprint: sshKey.fingerprint_sha_256,
                                        username: sshKey.username,
                                    })
                                }
                            />
                        ))}
                    </List>
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
                    placeholder="ssh-ed25519..., ssh-rsa..."
                    disabled={isCreatingKey}
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
                <APIError error={errorCreateKey} />
                {isCreatingKey && <Progress infinite />}
            </Popup>
        </Content>
    );
}
