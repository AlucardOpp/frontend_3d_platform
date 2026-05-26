export const FULL_DESCRIPTION_ACCEPT = '.txt,.doc,.docx,.pdf';

const ALLOWED_EXTENSIONS = ['.txt', '.doc', '.docx', '.pdf'];

export function isFullDescriptionFile(filename) {
  const index = filename.lastIndexOf('.');
  if (index === -1) {
    return false;
  }
  const extension = filename.slice(index).toLowerCase();
  return ALLOWED_EXTENSIONS.includes(extension);
}
