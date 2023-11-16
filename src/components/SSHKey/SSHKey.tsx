import {
    Button,
    List,
    ListActions,
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListProps,
    ListTitle,
    MaterialIcon,
} from "@vertex-center/components";

export function SSHKeys(props: Readonly<ListProps>) {
    return <List {...props} />;
}

type SSHKeyProps = {
    type: string;
    fingerprint: string;
    onDelete: () => void;
};

export default function SSHKey(props: Readonly<SSHKeyProps>) {
    const { type, fingerprint, onDelete } = props;

    return (
        <ListItem>
            <ListIcon>
                <MaterialIcon icon="key" />
            </ListIcon>
            <ListInfo>
                <ListTitle>SSH Key</ListTitle>
                <ListDescription>
                    {type} - {fingerprint}
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
