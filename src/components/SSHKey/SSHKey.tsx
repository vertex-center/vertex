import Symbol from "../Symbol/Symbol";
import ListItem from "../List/ListItem";
import ListSymbol from "../List/ListSymbol";
import ListInfo from "../List/ListInfo";
import ListTitle from "../List/ListTitle";
import ListDescription from "../List/ListDescription";
import List, { ListProps } from "../List/List";
import Button from "../Button/Button";
import ListActions from "../List/ListActions";

export function SSHKeys(props: ListProps) {
    return <List {...props} />;
}

type SSHKeyProps = {
    type: string;
    fingerprint: string;
    onDelete: () => void;
};

export default function SSHKey(props: SSHKeyProps) {
    const { type, fingerprint, onDelete } = props;

    return (
        <ListItem>
            <ListSymbol>
                <Symbol name="key" />
            </ListSymbol>
            <ListInfo>
                <ListTitle>SSH Key</ListTitle>
                <ListDescription>
                    {type} - {fingerprint}
                </ListDescription>
            </ListInfo>
            <ListActions>
                <Button rightSymbol="delete" onClick={onDelete}>
                    Delete
                </Button>
            </ListActions>
        </ListItem>
    );
}
