import Popup, { PopupActions } from "../../../../components/Popup/Popup";
import { Button, FormItem, Input, Vertical } from "@vertex-center/components";
import { APIError } from "../../../../components/Error/APIError";
import { useCreateContainer } from "../../hooks/useCreateContainer";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { DownloadSimple } from "@phosphor-icons/react";
import * as yup from "yup";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";

type Props = {
    dismiss: () => void;
};

const schema = yup
    .object({
        image: yup.string().required(),
    })
    .required();

export default function ManualInstallPopup(props: Readonly<Props>) {
    const { dismiss } = props;

    const {
        register,
        handleSubmit,
        watch,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(schema),
    });

    const { createContainer, isCreatingContainer, errorCreatingContainer } =
        useCreateContainer({
            onSuccess: dismiss,
        });

    const onSubmit = handleSubmit((data) => {
        const req = formattedImage(data.image);
        createContainer(req);
    });

    const formattedImage = (image?: string) => {
        if (!image) return;
        const fmt = {
            image: image,
            image_tag: "latest",
        };
        const split = image.split(":");
        if (split.length === 2) {
            fmt.image = split[0];
            fmt.image_tag = split[1];
        }
        return fmt;
    };

    const actions = (
        <PopupActions>
            <Button variant="outlined" onClick={dismiss}>
                Cancel
            </Button>
            <Button
                type="submit"
                variant="colored"
                rightIcon={<DownloadSimple />}
            >
                Install
            </Button>
        </PopupActions>
    );

    const image = watch("image");
    const formatted = formattedImage(image);
    let description = undefined;
    if (image) {
        description = `${formatted?.image}:${formatted?.image_tag}`;
    }

    return (
        <Popup onDismiss={dismiss} title="Install from Docker Registry">
            <form onSubmit={onSubmit}>
                <Vertical gap={20}>
                    <FormItem
                        label="Image"
                        error={errors.image?.message?.toString()}
                        description={description}
                        required
                    >
                        <Input
                            {...register("image")}
                            placeholder="postgres"
                            disabled={isCreatingContainer}
                        />
                    </FormItem>
                    {actions}
                </Vertical>
            </form>
            <ProgressOverlay show={isCreatingContainer} />
            <APIError error={errorCreatingContainer} />
        </Popup>
    );
}
