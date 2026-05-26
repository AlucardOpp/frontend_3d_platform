export const MODEL_3D_EXTENSIONS = [
  '.glb',
  '.gltf',
  '.obj',
  '.fbx',
  '.stl',
  '.dae',
  '.ply',
  '.3ds',
  '.wrl',
  '.vrml',
];

export const MODEL_ACCEPT_ATTRIBUTE = [
  'image/*',
  'video/*',
  ...MODEL_3D_EXTENSIONS,
].join(',');

export function getFileExtension(filename) {
  const index = filename.lastIndexOf('.');
  if (index === -1) {
    return '';
  }
  return filename.slice(index).toLowerCase();
}

export function isModel3DFile(filename) {
  return MODEL_3D_EXTENSIONS.includes(getFileExtension(filename));
}

export function isGlbFile(filename) {
  return getFileExtension(filename) === '.glb';
}

export function needsGlbConversion(filename) {
  return isModel3DFile(filename) && !isGlbFile(filename);
}
