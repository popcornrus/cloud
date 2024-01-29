export const Toast = async (config = {}) => {
    const {default: Toast} = await import('$lib/components/Toast.svelte');

    new Toast({
        target: document.body,
        props: {
            message: config?.message,
            type: config?.type,
            duration: config?.duration,
        }
    })
}