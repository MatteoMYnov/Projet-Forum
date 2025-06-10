# Fix : Images par défaut pour les anciens utilisateurs

## 🔍 **Problème identifié**

Les utilisateurs créés **avant** l'implémentation du système d'upload d'images avaient le champ `profile_picture` à `NULL` dans la base de données, ce qui causait l'affichage de "Photo de profil" au lieu de l'image par défaut.

## 🛠️ **Solution appliquée**

### 1. **Script de migration automatique**
```bash
go run scripts/migrate_default_avatars.go
```

Ce script :
- ✅ Détecte tous les utilisateurs avec `profile_picture = NULL` ou `profile_picture = ''`
- ✅ Les met à jour avec l'image par défaut `/img/avatars/default-avatar.png`
- ✅ Affiche un rapport de migration

### 2. **Logs de débogage ajoutés**
```go
if user.ProfilePicture != nil && *user.ProfilePicture != "" {
    profilePicture = *user.ProfilePicture
    log.Printf("🖼️ Utilisation image personnalisée: %s", profilePicture)
} else {
    log.Printf("🖼️ Utilisation image par défaut: %s (ProfilePicture=%v)", profilePicture, user.ProfilePicture)
}
```

## 📊 **Résultat de la migration**

```
🔄 Migration des avatars par défaut...
📊 2 utilisateurs sans photo de profil trouvés
✅ Migration terminée: 2 utilisateurs mis à jour avec l'avatar par défaut
```

### **Utilisateurs mis à jour :**
- `caca` (ID: 3) → maintenant avec avatar par défaut
- `caaca2` (ID: 5) → maintenant avec avatar par défaut

### **Utilisateurs déjà configurés :**
- `romain` (ID: 6) → garde son image personnalisée

## ✅ **Test de vérification**

1. **Connectez-vous comme "caca"**
2. **Allez sur `/profile`**
3. **L'image par défaut devrait maintenant s'afficher**

## 🔧 **Prévention future**

Pour tous les **nouveaux utilisateurs**, l'image par défaut est automatiquement assignée lors de l'inscription :

```go
// Dans RegisterHandler
defaultPath := c.uploadService.GetDefaultAvatarPath()
profilePicturePath = &defaultPath
```

## 📝 **Commandes utiles**

### Vérifier les utilisateurs sans avatar :
```sql
SELECT id_user, username, profile_picture 
FROM users 
WHERE profile_picture IS NULL OR profile_picture = '';
```

### Vérifier que la migration a fonctionné :
```sql
SELECT id_user, username, profile_picture 
FROM users 
ORDER BY id_user;
```

### Relancer la migration si nécessaire :
```bash
go run scripts/migrate_default_avatars.go
```

## 🎯 **État final**

- ✅ Tous les utilisateurs ont maintenant une image de profil (par défaut ou personnalisée)
- ✅ Les nouveaux utilisateurs reçoivent automatiquement l'image par défaut
- ✅ Le système fonctionne de manière cohérente pour tous les cas 