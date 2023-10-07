import Icon from "../Icon/Icon";
import ListItem from "../List/ListItem";
import ListIcon from "../List/ListIcon";
import ListInfo from "../List/ListInfo";
import ListTitle from "../List/ListTitle";
import ListDescription from "../List/ListDescription";
import List, { ListProps } from "../List/List";
import Button from "../Button/Button";
import ListActions from "../List/ListActions";

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
                <Icon name="key" />
            </ListIcon>
            <ListInfo>
                <ListTitle>SSH Key</ListTitle>
                <ListDescription>
                    {type} - {fingerprint}
                </ListDescription>
            </ListInfo>
            <ListActions>
                <Button rightIcon="delete" onClick={onDelete}>
                    Delete
                </Button>
            </ListActions>
        </ListItem>
    );
}
