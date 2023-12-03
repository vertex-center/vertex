import Content from "../../../../components/Content/Content";
import {
    Button,
    Code,
    List,
    ListActions,
    ListDescription,
    ListInfo,
    ListItem,
    ListTitle,
    MaterialIcon,
    Paragraph,
    TextField,
    Title,
} from "@vertex-center/components";
import {
    useCreateCurrentUserEmail,
    useCurrentUserEmails,
    useDeleteCurrentUserEmail,
} from "../../hooks/useEmails";
import Progress, {
    ProgressOverlay,
} from "../../../../components/Progress/Progress";
import { APIError } from "../../../../components/Error/APIError";
import NoItems from "../../../../components/NoItems/NoItems";
import Popup from "../../../../components/Popup/Popup";
import { ChangeEvent, Fragment, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";

export default function AccountEmails() {
    const queryClient = useQueryClient();

    const [showCreatePopup, setShowCreatePopup] = useState(false);
    const [showDeletePopup, setShowDeletePopup] = useState(false);

    const [email, setEmail] = useState("");

    const { emails, isLoadingEmails, errorEmails } = useCurrentUserEmails();
    const { createEmail, isCreatingEmail, errorCreateEmail, resetCreateEmail } =
        useCreateCurrentUserEmail({
            onSuccess: () => {
                queryClient.invalidateQueries(["user_emails"]);
                dismissCreatePopup();
            },
        });
    const { deleteEmail, isDeletingEmail, errorDeleteEmail, resetDeleteEmail } =
        useDeleteCurrentUserEmail({
            onSuccess: () => {
                queryClient.invalidateQueries(["user_emails"]);
                dismissDeletePopup();
            },
        });

    const dismissCreatePopup = () => {
        setShowCreatePopup(false);
        setEmail("");
        resetCreateEmail();
    };

    const dismissDeletePopup = () => {
        setShowDeletePopup(false);
    };

    const openDeletePopup = (email: string) => {
        setEmail(email);
        setShowDeletePopup(true);
    };

    const onEmailChange = (e: ChangeEvent<HTMLInputElement>) => {
        setEmail(e.target.value);
    };

    const popupCreateActions = (
        <Fragment>
            <Button variant="outlined" onClick={dismissCreatePopup}>
                Cancel
            </Button>
            <Button
                variant="colored"
                onClick={() => createEmail({ email })}
                rightIcon={<MaterialIcon icon="add" />}
            >
                Add
            </Button>
        </Fragment>
    );

    const popupDeleteActions = (
        <Fragment>
            <Button variant="outlined" onClick={dismissDeletePopup}>
                Cancel
            </Button>
            <Button
                variant="danger"
                onClick={() => deleteEmail({ email })}
                rightIcon={<MaterialIcon icon="delete" />}
            >
                Delete
            </Button>
        </Fragment>
    );

    return (
        <Content>
            <Title variant="h2">Emails</Title>
            <ProgressOverlay show={isLoadingEmails} />
            <APIError error={errorEmails} />
            {!isLoadingEmails && !emails?.length ? (
                <NoItems
                    text="You don't have any email address yet."
                    icon="email"
                />
            ) : (
                <List>
                    {emails?.map((m) => (
                        <ListItem key={m.id}>
                            <ListInfo>
                                <ListTitle>{m.email}</ListTitle>
                                <ListDescription>
                                    Added on{" "}
                                    {new Date(
                                        m.created_at * 1000
                                    ).toDateString()}
                                </ListDescription>
                            </ListInfo>
                            <ListActions>
                                <Button
                                    variant="danger"
                                    rightIcon={<MaterialIcon icon="delete" />}
                                    onClick={() => openDeletePopup(m.email)}
                                >
                                    Delete
                                </Button>
                            </ListActions>
                        </ListItem>
                    ))}
                </List>
            )}
            <div>
                <Button
                    variant="colored"
                    leftIcon={<MaterialIcon icon="add" />}
                    disabled={isLoadingEmails}
                    onClick={() => setShowCreatePopup(true)}
                >
                    Add email
                </Button>
            </div>
            <Popup
                show={showCreatePopup}
                onDismiss={dismissCreatePopup}
                title="Add email address"
                actions={popupCreateActions}
            >
                <TextField
                    label="Email address"
                    type="email"
                    value={email}
                    onChange={onEmailChange}
                    required
                />
                <APIError error={errorCreateEmail} />
                {isCreatingEmail && <Progress />}
            </Popup>
            <Popup
                show={showDeletePopup}
                onDismiss={dismissDeletePopup}
                title="Delete email?"
                actions={popupDeleteActions}
            >
                <Paragraph>
                    Are you sure you want to delete this email address?
                </Paragraph>
                <Code language={"text"}>{email}</Code>
                <APIError error={errorDeleteEmail} />
                {isDeletingEmail && <Progress />}
            </Popup>
        </Content>
    );
}
