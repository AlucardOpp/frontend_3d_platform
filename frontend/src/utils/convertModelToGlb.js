import * as THREE from 'three';
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js';
import { GLTFExporter } from 'three/examples/jsm/exporters/GLTFExporter.js';
import { DRACOLoader } from 'three/examples/jsm/loaders/DRACOLoader.js';
import { OBJLoader } from 'three/examples/jsm/loaders/OBJLoader.js';
import { FBXLoader } from 'three/examples/jsm/loaders/FBXLoader.js';
import { STLLoader } from 'three/examples/jsm/loaders/STLLoader.js';
import { ColladaLoader } from 'three/examples/jsm/loaders/ColladaLoader.js';
import { PLYLoader } from 'three/examples/jsm/loaders/PLYLoader.js';
import { TDSLoader } from 'three/examples/jsm/loaders/TDSLoader.js';
import { VRMLLoader } from 'three/examples/jsm/loaders/VRMLLoader.js';
import { getFileExtension } from './model3dFormats';

const DRACO_DECODER_PATH = 'https://www.gstatic.com/draco/versioned/decoders/1.5.6/';

function createGLTFLoader() {
  const loader = new GLTFLoader();
  const dracoLoader = new DRACOLoader();
  dracoLoader.setDecoderPath(DRACO_DECODER_PATH);
  loader.setDRACOLoader(dracoLoader);
  return loader;
}

function wrapGeometryAsScene(geometry) {
  const material = new THREE.MeshStandardMaterial({ color: 0xcccccc });
  const mesh = new THREE.Mesh(geometry, material);
  geometry.computeVertexNormals();
  const scene = new THREE.Scene();
  scene.add(mesh);
  return scene;
}

function toScene(object) {
  if (object.isScene) {
    return object;
  }
  const scene = new THREE.Scene();
  scene.add(object);
  return scene;
}

async function loadModelScene(file) {
  const extension = getFileExtension(file.name);
  const url = URL.createObjectURL(file);

  try {
    switch (extension) {
      case '.glb':
      case '.gltf': {
        const gltf = await createGLTFLoader().loadAsync(url);
        return toScene(gltf.scene);
      }
      case '.obj': {
        const object = await new OBJLoader().loadAsync(url);
        return toScene(object);
      }
      case '.fbx': {
        const object = await new FBXLoader().loadAsync(url);
        return toScene(object);
      }
      case '.stl': {
        const geometry = await new STLLoader().loadAsync(url);
        return wrapGeometryAsScene(geometry);
      }
      case '.ply': {
        const geometry = await new PLYLoader().loadAsync(url);
        return wrapGeometryAsScene(geometry);
      }
      case '.dae': {
        const collada = await new ColladaLoader().loadAsync(url);
        return toScene(collada.scene);
      }
      case '.3ds': {
        const object = await new TDSLoader().loadAsync(url);
        return toScene(object);
      }
      case '.wrl':
      case '.vrml': {
        const object = await new VRMLLoader().loadAsync(url);
        return toScene(object);
      }
      default:
        throw new Error(`Формат ${extension} не поддерживается для конвертации`);
    }
  } finally {
    URL.revokeObjectURL(url);
  }
}

function exportSceneToGlb(scene) {
  return new Promise((resolve, reject) => {
    const exporter = new GLTFExporter();
    exporter.parse(
      scene,
      (result) => {
        if (!(result instanceof ArrayBuffer)) {
          reject(new Error('Ожидался бинарный GLB'));
          return;
        }
        resolve(result);
      },
      (error) => reject(error),
      { binary: true },
    );
  });
}

function createGlbFile(originalFile, glbBuffer) {
  const baseName = originalFile.name.replace(/\.[^.]+$/, '') || 'model';
  return new File([glbBuffer], `${baseName}.glb`, {
    type: 'model/gltf-binary',
    lastModified: Date.now(),
  });
}

export async function convertModelToGlb(file) {
  const scene = await loadModelScene(file);
  const glbBuffer = await exportSceneToGlb(scene);
  return createGlbFile(file, glbBuffer);
}
