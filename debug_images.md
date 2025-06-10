# Débogage des images de profil

## Problème résolu ✅

### 🔍 **Problème identifié**
- Le template de profil utilisait des chemins relatifs incorrects
- Mauvais dossier de référence (`/img/avatar/` au lieu de `/img/avatars/`)

### 🛠️ **Solutions appliquées**

1. **Correction des chemins dans le processTemplate**
   ```go
   // Remplacement de ../img/avatar/photo-profil.jpg par le vrai chemin
   profilePicture := "/img/avatars/default-avatar.png"
   if user.ProfilePicture != nil && *user.ProfilePicture != "" {
       profilePicture = *user.ProfilePicture
   }
   ```

2. **Mise à jour du template processing**
   - Remplace `src="../img/avatar/photo-profil.jpg"` ✅
   - Remplace `src="../img/avatar/avatar-utilisateur.jpg"` ✅

### 🧪 **Tests à effectuer**

1. **Test de l'image par défaut**
   - URL: `http://localhost:2557/img/avatars/default-avatar.png`
   - Doit afficher l'image par défaut

2. **Test d'une image uploadée**
   - URL: `http://localhost:2557/img/avatars/1749581819_7363194bd43016225c90b1d43a85778b.png`
   - Doit afficher l'image de l'utilisateur "romain"

3. **Test du profil**
   - Aller sur `/profile` connecté comme "romain"
   - L'image de profil doit maintenant s'afficher

### 📂 **Structure des URLs**
```
Serveur de fichiers statiques:
http://localhost:2557/img/avatars/filename.ext

Mapping dans main.go:
/img/ → ./website/img/

Donc:
/img/avatars/default-avatar.png → ./website/img/avatars/default-avatar.png
```

### 🔄 **Chemin de remplacement dans le template**

**Avant:**
```html
<img src="../img/avatar/photo-profil.jpg" alt="Photo de profil" class="profile-pic" />
```

**Après:**
```html
<img src="/img/avatars/default-avatar.png" alt="Photo de profil" class="profile-pic" />
<!-- OU -->
<img src="/img/avatars/1749581819_7363194bd43016225c90b1d43a85778b.png" alt="Photo de profil" class="profile-pic" />
```

### ✅ **Vérification rapide**
1. Connectez-vous comme "romain" 
2. Allez sur `/profile`
3. L'image doit maintenant être visible
4. Inspectez l'élément pour vérifier l'URL de l'image 