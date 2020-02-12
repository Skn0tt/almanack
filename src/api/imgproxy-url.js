export default function imageURL(filepath, { width = 400, height = 300 } = {}) {
  if (!filepath) {
    return "";
  }
  let baseURL = "https://images.data.spotlightpa.org";
  let signature = "insecure";
  let resizing_type = "auto";
  let gravity = "sm";
  let enlarge = "1";
  let quality = "75";
  let encoded_source_url = btoa(filepath);
  let extension = "jpeg";

  return `${baseURL}/${signature}/rs:${resizing_type}:${width}:${height}/g:${gravity}/el:${enlarge}/q:${quality}/${encoded_source_url}.${extension}`;
}
