import { isModel3DFile, isGlbFile } from './model3dFormats';

export async function prepareModelFileForUpload(file) {
  if (!isModel3DFile(file.name) || isGlbFile(file.name)) {
    return file;
  }

  const { convertModelToGlb } = await import(
    /* webpackChunkName: "model-converter" */ './convertModelToGlb'
  );
  return convertModelToGlb(file);
}
