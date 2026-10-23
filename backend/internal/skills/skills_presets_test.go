package skills

import "testing"

func TestAudioBuiltinSkillCategoryIsPublished(t *testing.T) {
	if got := skillCategoryLabels["audio"]; got != "音频创作" {
		t.Fatalf("audio category label = %q, want 音频创作", got)
	}
	packages, err := loadBuiltinSkillPackages(nil)
	if err != nil {
		t.Fatal(err)
	}
	foundAudioSkill := false
	for _, item := range packages {
		if item.skill.ID == "doubao-audio-prompt" && item.skill.Tag == "creative" {
			foundAudioSkill = true
			break
		}
	}
	if !foundAudioSkill {
		t.Fatal("doubao-audio-prompt builtin must retain the Dev creative category")
	}
	for _, category := range skillCategories() {
		if category.Value == "audio" && category.Label == "音频创作" {
			return
		}
	}
	t.Fatal("audio category missing from skills API catalog")
}

func TestSkillPresetsCatalogIsValid(t *testing.T) {
	presets, err := New(nil, "", nil).SkillPresets()
	if err != nil {
		t.Fatal(err)
	}
	if len(presets) == 0 {
		t.Fatal("场景预设目录不能为空")
	}
	seeded := make(map[string]struct{}, 128)
	ids, err := builtinSeedSkillIDs()
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range ids {
		seeded[id] = struct{}{}
	}
	for _, preset := range presets {
		if len(preset.SkillIDs) == 0 || len(preset.SkillIDs) > 8 {
			t.Fatalf("%s 技能数超出 1-8", preset.PresetID)
		}
		for _, skillID := range preset.SkillIDs {
			if _, ok := seeded[skillID]; !ok {
				t.Fatalf("%s 引用未上架技能 %s", preset.PresetID, skillID)
			}
		}
	}
}
