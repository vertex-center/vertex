import "./ProfilePicture.sass";
import cx from "classnames";
import { HTMLProps } from "react";

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
        return <div style={{ width: size, height: size }} {...properties} />;
    }
    return <img alt={alt} width={size} height={size} {...properties} />;
}
