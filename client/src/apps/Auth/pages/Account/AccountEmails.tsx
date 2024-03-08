import Content from "../../../../components/Content/Content";
import {
    Button,
    Code,
    FormItem,
    Input,
    List,
    ListActions,
    ListDescription,
    ListInfo,
    ListItem,
    ListTitle,
    Paragraph,
    Title,
    Vertical,
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
import Popup, { PopupActions } from "../../../../components/Popup/Popup";
import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Envelope, Plus, Trash } from "@phosphor-icons/react";
import { useForm } from "react-hook-form";
import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup";
import { Email } from "../../backend/models";

type DeleteEmailPopupProps = {
    show: boolean;
    onDismiss: () => void;
    email: Partial<Email>;
};

function DeleteEmailPopup(props: DeleteEmailPopupProps) {
    const queryClient = useQueryClient();

    const { show, onDismiss, email } = props;

    const { deleteEmail, isDeletingEmail, errorDeleteEmail } =
        useDeleteCurrentUserEmail({
            onSuccess: async () => {
                await queryClient.invalidateQueries({
                    queryKey: ["user_emails"],
                });
                onDismiss();
            },
        });

    return (
        <Popup show={show} onDismiss={onDismiss} title="Delete email?">
            <Paragraph>
                Are you sure you want to delete this email address?
            </Paragraph>
            <Code language={"text"}>{email.email}</Code>
            <APIError error={errorDeleteEmail} />
            {isDeletingEmail && <Progress />}
            <PopupActions>
                <Button variant="outlined" onClick={onDismiss}>
                    Cancel
                </Button>
                <Button
                    variant="danger"
                    rightIcon={<Trash />}
                    onClick={() => deleteEmail(email)}
                >
                    Delete
                </Button>
            </PopupActions>
        </Popup>
    );
}

type CreateEmailPopupProps = {
    show: boolean;
    onDismiss: () => void;
};

const createEmailSchema = yup.object().shape({
    email: yup.string().email().required(),
});

function CreateEmailPopup(props: CreateEmailPopupProps) {
    const queryClient = useQueryClient();

    const { show, onDismiss: _onDismiss } = props;

    const {
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(createEmailSchema),
    });

    const onDismiss = () => {
        reset();
        _onDismiss();
    };

    const { createEmail, isCreatingEmail, errorCreateEmail } =
        useCreateCurrentUserEmail({
            onSuccess: async () => {
                await queryClient.invalidateQueries({
                    queryKey: ["user_emails"],
                });
                onDismiss();
                reset();
            },
        });

    const onSubmit = handleSubmit((data) => {
        createEmail(data);
    });

    return (
        <Popup show={show} onDismiss={onDismiss} title="Add email address">
            <form onSubmit={onSubmit}>
                <Vertical gap={12}>
                    <FormItem
                        label="Email address"
                        error={errors.email?.message?.toString()}
                        required
                    >
                        <Input type="email" {...register("email")} />
                    </FormItem>
                    <APIError error={errorCreateEmail} />
                    {isCreatingEmail && <Progress />}
                    <PopupActions>
                        <Button variant="outlined" onClick={onDismiss}>
                            Cancel
                        </Button>
                        <Button
                            type="submit"
                            variant="colored"
                            rightIcon={<Plus />}
                        >
                            Add
                        </Button>
                    </PopupActions>
                </Vertical>
            </form>
        </Popup>
    );
}

export default function AccountEmails() {
    const [showCreatePopup, setShowCreatePopup] = useState(false);
    const [showDeletePopup, setShowDeletePopup] = useState(false);

    const [email, setEmail] = useState<Partial<Email>>({});

    const { emails, isLoadingEmails, errorEmails } = useCurrentUserEmails();

    const dismissCreatePopup = () => setShowCreatePopup(false);
    const dismissDeletePopup = () => setShowDeletePopup(false);

    const openDeletePopup = (email: Partial<Email>) => {
        setEmail(email);
        setShowDeletePopup(true);
    };

    return (
        <Content>
            <Title variant="h2">Emails</Title>
            <ProgressOverlay show={isLoadingEmails} />
            <APIError error={errorEmails} />
            {!isLoadingEmails && !emails?.length ? (
                <NoItems
                    text="You don't have any email address yet."
                    icon={<Envelope />}
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
                                    rightIcon={<Trash />}
                                    onClick={() => openDeletePopup(m)}
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
                    leftIcon={<Plus />}
                    disabled={isLoadingEmails}
                    onClick={() => setShowCreatePopup(true)}
                >
                    Add email
                </Button>
            </div>
            <CreateEmailPopup
                show={showCreatePopup}
                onDismiss={dismissCreatePopup}
            />
            <DeleteEmailPopup
                show={showDeletePopup}
                onDismiss={dismissDeletePopup}
                email={email}
            />
        </Content>
    );
}
