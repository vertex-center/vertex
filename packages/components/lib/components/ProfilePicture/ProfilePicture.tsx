import "./ProfilePicture.sass";
import cx from "classnames";
import { HTMLProps } from "react";
import { User } from "@phosphor-icons/react";

export type ProfilePictureProps = HTMLProps<HTMLImageElement> & {
    size?: number;
};

export function ProfilePicture(props: Readonly<ProfilePictureProps>) {
    const { className, alt, size = 40, ...others } = props;

    const properties = {
        className: cx("profile-picture", className),
        ...others,
    };

    if (props.src === undefined) {
        return (
            <div style={{ width: size, height: size }} {...properties}>
                <User size={20} />
            </div>
        );
    }
    return <img alt={alt} width={size} height={size} {...properties} />;
}
