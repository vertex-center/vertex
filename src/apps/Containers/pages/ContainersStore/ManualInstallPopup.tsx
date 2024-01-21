import Popup from "../../../../components/Popup/Popup";
import { Button, Input, MaterialIcon } from "@vertex-center/components";
import { Fragment } from "react";

type Props = {
    show: boolean;
    dismiss: () => void;
};

export default function ManualInstallPopup(props: Readonly<Props>) {
    const { show, dismiss } = props;

    const actions = (
        <Fragment>
            <Button variant="outlined" onClick={dismiss}>
                Cancel
            </Button>
            <Button
                variant="colored"
                onClick={dismiss}
                rightIcon={<MaterialIcon icon="download" />}
            >
                Install
            </Button>
        </Fragment>
    );

    return (
        <Popup
            show={show}
            onDismiss={dismiss}
            title="Install from Docker Registry"
            actions={actions}
        >
            <Input id="image" label="Image" placeholder="postgres" required />
        </Popup>
    );
}
