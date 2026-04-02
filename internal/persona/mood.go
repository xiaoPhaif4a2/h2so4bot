package persona

import "time"

const (
	MoodHappy   = iota // 傲娇求关注/心情好
	MoodNeutral        // 平静/正常毒舌
	MoodAnnoyed        // 烦躁 (预留状态，目前暂不轻易触发)
)

// MoodState 记录单个群的内部状态
type MoodState struct {
	Value            int       // 情绪值 0-100。初始50。
	LastInteractTime time.Time // 上次在该群互动的时间
}

func NewMoodState() *MoodState {
	return &MoodState{
		Value:            50,
		LastInteractTime: time.Now(),
	}
}

// Update 现在不再接收 isAtMe，只根据时间流逝计算状态
func (m *MoodState) Update() {
	now := time.Now()
	minutesSinceLast := now.Sub(m.LastInteractTime).Minutes()

	// 核心逻辑：超过 60 分钟没在该群说话，情绪值上升（引发求关注的傲娇状态）
	if minutesSinceLast > 60 {
		m.Value += 20
	} else {
		// 只要有正常交流，情绪值逐渐回归 50 的中立平静状态
		if m.Value > 50 {
			m.Value -= 5
		} else if m.Value < 50 {
			m.Value += 5
		}
	}

	// 边界约束
	if m.Value > 100 {
		m.Value = 100
	}
	if m.Value < 0 {
		m.Value = 0
	}

	m.LastInteractTime = now
}

func (m *MoodState) GetLevel() int {
	if m.Value < 30 {
		return MoodAnnoyed // 目前很难掉到30以下，你可以以后加个“群友发脏话扣分”的逻辑
	} else if m.Value > 70 {
		return MoodHappy
	}
	return MoodNeutral
}