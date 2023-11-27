import {
    Button,
    ListActions,
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
    MaterialIcon,
} from "@vertex-center/components";

type SSHKeyProps = {
    type: string;
    fingerprint: string;
    username: string;
    onDelete: () => void;
};

export default function SSHKey(props: Readonly<SSHKeyProps>) {
    const { type, fingerprint, username, onDelete } = props;

    return (
        <ListItem>
            <ListIcon>
                <MaterialIcon icon="key" />
            </ListIcon>
            <ListInfo>
                <ListTitle>SSH Key</ListTitle>
                <ListDescription>
                    {type} - {fingerprint} - {username}
                </ListDescription>
            </ListInfo>
            <ListActions>
                <Button
                    variant="danger"
                    rightIcon={<MaterialIcon icon="delete" />}
                    onClick={onDelete}
                >
                    Delete
                </Button>
            </ListActions>
        </ListItem>
    );
}
