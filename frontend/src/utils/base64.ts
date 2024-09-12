export async function ToBase64(file): Promise<any> {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result?.toString().replace(/^data:image\/[a-z]+;base64,/, ""));
        reader.onerror = error => reject(error);
    });
}