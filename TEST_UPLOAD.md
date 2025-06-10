# Test de la fonctionnalité d'upload d'image de profil

## Fonctionnalités implémentées ✅

### 1. **Page d'inscription avec upload d'image**
- ✅ Formulaire avec section photo de profil
- ✅ Aperçu en temps réel de l'image sélectionnée
- ✅ Validation côté client (taille max 5MB, types d'images)
- ✅ Bouton pour supprimer l'image sélectionnée
- ✅ Image par défaut si aucune image n'est fournie

### 2. **Backend upload et gestion**
- ✅ Service d'upload sécurisé (`UploadService`)
- ✅ Génération de noms de fichiers uniques
- ✅ Validation des types de fichiers et tailles
- ✅ Sauvegarde dans `website/img/avatars/`
- ✅ Nettoyage automatique en cas d'erreur

### 3. **Base de données et modèles**
- ✅ Champ `profile_picture` dans le modèle User
- ✅ Support dans RegisterRequest
- ✅ Mise à jour des requêtes SQL (CREATE et SELECT)

### 4. **Affichage des avatars**
- ✅ Avatars utilisateur dans la liste des threads
- ✅ Avatar utilisateur dans la page de profil
- ✅ Fallback vers l'image par défaut si aucune image

## Comment tester

### 1. **Test d'inscription sans image**
1. Aller sur `localhost:2557/register`
2. Remplir le formulaire sans sélectionner d'image
3. S'inscrire
4. Vérifier que l'utilisateur a l'image par défaut

### 2. **Test d'inscription avec image**
1. Aller sur `localhost:2557/register`
2. Cliquer sur la zone d'avatar pour sélectionner une image
3. Choisir une image JPG/PNG (< 5MB)
4. Vérifier l'aperçu en temps réel
5. Remplir le reste du formulaire et s'inscrire
6. Se connecter et vérifier que l'image apparaît dans le profil

### 3. **Test de validation**
1. Essayer avec un fichier > 5MB → Erreur
2. Essayer avec un fichier non-image → Erreur
3. Vérifier que les erreurs s'affichent correctement

### 4. **Test d'affichage**
1. Créer des threads avec des utilisateurs ayant différents avatars
2. Vérifier que les avatars s'affichent correctement dans `/threads`
3. Vérifier l'avatar dans `/profile`

## Structure des fichiers

```
website/img/avatars/
├── default-avatar.png          # Image par défaut
├── 1703123456_abc123def.jpg   # Images uploadées (timestamp_unique.ext)
└── 1703123789_xyz789abc.png   # Format: timestamp_uniqueID.extension
```

## Images supportées
- **Types**: JPG, JPEG, PNG, GIF, WEBP
- **Taille max**: 5MB
- **Validation**: Côté client ET serveur

## Sécurité
- ✅ Validation stricte des types de fichiers
- ✅ Limite de taille appliquée
- ✅ Noms de fichiers uniques (évite les conflits)
- ✅ Nettoyage automatique en cas d'erreur
- ✅ Pas d'exécution de code uploadé

## Notes techniques
- Les images sont stockées localement dans `website/img/avatars/`
- Le chemin est sauvegardé en base comme `/img/avatars/filename.ext`
- Le service d'upload est réutilisable pour d'autres types de fichiers
- Support du format multipart/form-data pour les uploads 