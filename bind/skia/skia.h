
typedef unsigned char      uint8_t;   // 无符号8位整数
typedef unsigned short     uint16_t;  // 无符号16位整数
typedef unsigned int       uint32_t;  // 无符号32位整数
typedef unsigned long long uint64_t;  // 无符号64位整数
typedef signed char        int8_t;    // 有符号8位整数
typedef signed short       int16_t;   // 有符号16位整数
typedef signed int         int32_t;   // 有符号32位整数
typedef signed long long   int64_t;   // 有符号64位整数
typedef unsigned char bool;           // 使用 typedef 定义 bool 类型

typedef int* intptr_t;

typedef unsigned char      uint8_t;   // 无符号8位整数
typedef unsigned short     uint16_t;  // 无符号16位整数
typedef unsigned int       uint32_t;  // 无符号32位整数
typedef unsigned long long uint64_t;  // 无符号64位整数
typedef signed char        int8_t;    // 有符号8位整数
typedef signed short       int16_t;   // 有符号16位整数
typedef signed int         int32_t;   // 有符号32位整数
typedef signed long long   int64_t;   // 有符号64位整数
typedef unsigned char bool;           // 使用 typedef 定义 bool 类型
//c/sk_types.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_types_DEFINED
#define sk_types_DEFINED

////#include <stdint.h>
////#include <stddef.h>

#ifdef __cplusplus
    #define SK_C_PLUS_PLUS_BEGIN_GUARD    extern "C" {
    #define SK_C_PLUS_PLUS_END_GUARD      }
#else
    ////#include <stdbool.h>
    #define SK_C_PLUS_PLUS_BEGIN_GUARD
    #define SK_C_PLUS_PLUS_END_GUARD
#endif

#if !defined(SK_C_API)
    #if defined(SKIA_C_DLL)
        #if defined(_MSC_VER)
            #if SKIA_IMPLEMENTATION
                #define SK_C_API __declspec(dllexport)
            #else
                #define SK_C_API __declspec(dllimport)
            #endif
        #else
            #define SK_C_API __attribute__((visibility("default")))
        #endif
    #else
        #define SK_C_API
    #endif
#endif

#if defined(_WIN32)
    // On Windows, Vulkan commands use the stdcall convention
    #define VKAPI_ATTR
    #define VKAPI_CALL __stdcall
    #define VKAPI_PTR  VKAPI_CALL
#elif defined(__ANDROID__) && defined(__ARM_ARCH) && __ARM_ARCH < 7
    #error "Vulkan isn't supported for the 'armeabi' NDK ABI"
#elif defined(__ANDROID__) && defined(__ARM_ARCH) && __ARM_ARCH >= 7 && defined(__ARM_32BIT_STATE)
    // On Android 32-bit ARM targets, Vulkan functions use the "hardfloat"
    // calling convention, i.e. float parameters are passed in registers. This
    // is true even if the rest of the application passes floats on the stack,
    // as it does by default when compiling for the armeabi-v7a NDK ABI.
    #define VKAPI_ATTR __attribute__((pcs("aapcs-vfp")))
    #define VKAPI_CALL
    #define VKAPI_PTR  VKAPI_ATTR
#else
    // On other platforms, use the default calling convention
    #define VKAPI_ATTR
    #define VKAPI_CALL
    #define VKAPI_PTR
#endif

#if !defined(SK_TO_STRING)
    #define SK_TO_STRING(X) SK_TO_STRING_IMPL(X)
    #define SK_TO_STRING_IMPL(X) #X
#endif

#ifndef SK_C_INCREMENT
#define SK_C_INCREMENT 0
#endif

///////////////////////////////////////////////////////////////////////////////////////

SK_C_PLUS_PLUS_BEGIN_GUARD

typedef struct sk_refcnt_t sk_refcnt_t;
typedef struct sk_nvrefcnt_t sk_nvrefcnt_t;

typedef struct sk_flattenable_t sk_flattenable_t;

typedef uint32_t sk_color_t;
typedef uint32_t sk_pmcolor_t;

/* This macro assumes all arguments are >=0 and <=255. */
#define sk_color_set_argb(a, r, g, b)   (((a) << 24) | ((r) << 16) | ((g) << 8) | (b))
#define sk_color_get_a(c)               (((c) >> 24) & 0xFF)
#define sk_color_get_r(c)               (((c) >> 16) & 0xFF)
#define sk_color_get_g(c)               (((c) >>  8) & 0xFF)
#define sk_color_get_b(c)               (((c) >>  0) & 0xFF)

typedef struct sk_color4f_t_ {
    float fR;
    float fG;
    float fB;
    float fA;
} sk_color4f_t;

typedef enum sk_colortype_t_ {
    UNKNOWN_SK_COLORTYPE = 0,
    ALPHA_8_SK_COLORTYPE,
    RGB_565_SK_COLORTYPE,
    ARGB_4444_SK_COLORTYPE,
    RGBA_8888_SK_COLORTYPE,
    RGB_888X_SK_COLORTYPE,
    BGRA_8888_SK_COLORTYPE,
    RGBA_1010102_SK_COLORTYPE,
    BGRA_1010102_SK_COLORTYPE,
    RGB_101010X_SK_COLORTYPE,
    BGR_101010X_SK_COLORTYPE,
    BGR_101010X_XR_SK_COLORTYPE,
    GRAY_8_SK_COLORTYPE,
    RGBA_F16_NORM_SK_COLORTYPE,
    RGBA_F16_SK_COLORTYPE,
    RGBA_F32_SK_COLORTYPE,

    // READONLY
    R8G8_UNORM_SK_COLORTYPE,
    A16_FLOAT_SK_COLORTYPE,
    R16G16_FLOAT_SK_COLORTYPE,
    A16_UNORM_SK_COLORTYPE,
    R16G16_UNORM_SK_COLORTYPE,
    R16G16B16A16_UNORM_SK_COLORTYPE,
    SRGBA_8888_SK_COLORTYPE,
    R8_UNORM_SK_COLORTYPE,
} sk_colortype_t;

typedef enum sk_alphatype_t_ {
    UNKNOWN_SK_ALPHATYPE,
    OPAQUE_SK_ALPHATYPE,
    PREMUL_SK_ALPHATYPE,
    UNPREMUL_SK_ALPHATYPE,
} sk_alphatype_t;

typedef enum sk_pixelgeometry_t_ {
    UNKNOWN_SK_PIXELGEOMETRY,
    RGB_H_SK_PIXELGEOMETRY,
    BGR_H_SK_PIXELGEOMETRY,
    RGB_V_SK_PIXELGEOMETRY,
    BGR_V_SK_PIXELGEOMETRY,
} sk_pixelgeometry_t;

typedef enum sk_surfaceprops_flags_t_ {
    NONE_SK_SURFACE_PROPS_FLAGS = 0,
    USE_DEVICE_INDEPENDENT_FONTS_SK_SURFACE_PROPS_FLAGS = 1 << 0,
} sk_surfaceprops_flags_t;

typedef struct sk_surfaceprops_t sk_surfaceprops_t;

typedef struct sk_point_t_ {
    float   x;
    float   y;
} sk_point_t;

typedef sk_point_t sk_vector_t;

typedef struct sk_irect_t_ {
    int32_t left;
    int32_t top;
    int32_t right;
    int32_t bottom;
} sk_irect_t;

typedef struct sk_rect_t_ {
    float   left;
    float   top;
    float   right;
    float   bottom;
} sk_rect_t;

typedef struct sk_matrix_t_ {
    float scaleX,  skewX, transX;
    float  skewY, scaleY, transY;
    float persp0, persp1, persp2;
} sk_matrix_t;

// row major
typedef struct sk_matrix44_t_ {
    // name: m<row><col>
    float m00, m01, m02, m03; // row 0
    float m10, m11, m12, m13; // row 1
    float m20, m21, m22, m23; // row 2
    float m30, m31, m32, m33; // row 3
} sk_matrix44_t;

/**
    A sk_canvas_t encapsulates all of the state about drawing into a
    destination This includes a reference to the destination itself,
    and a stack of matrix/clip values.
*/
typedef struct sk_canvas_t sk_canvas_t;
typedef struct sk_nodraw_canvas_t sk_nodraw_canvas_t;
typedef struct sk_nway_canvas_t sk_nway_canvas_t;
typedef struct sk_overdraw_canvas_t sk_overdraw_canvas_t;
/**
    A sk_data_ holds an immutable data buffer.
*/
typedef struct sk_data_t sk_data_t;
/**
    A sk_drawable_t is a abstraction for drawings that changed while
    drawing.
*/
typedef struct sk_drawable_t sk_drawable_t;
/**
    A sk_image_t is an abstraction for drawing a rectagle of pixels.
    The content of the image is always immutable, though the actual
    storage may change, if for example that image can be re-created via
    encoded data or other means.
*/
typedef struct sk_image_t sk_image_t;
/**
    A sk_maskfilter_t is an object that perform transformations on an
    alpha-channel mask before drawing it; it may be installed into a
    sk_paint_t.  Each time a primitive is drawn, it is first
    scan-converted into a alpha mask, which os handed to the
    maskfilter, which may create a new mask is to render into the
    destination.
 */
typedef struct sk_maskfilter_t sk_maskfilter_t;
/**
    A sk_paint_t holds the style and color information about how to
    draw geometries, text and bitmaps.
*/
typedef struct sk_paint_t sk_paint_t;
typedef struct sk_font_t sk_font_t;
/**
    A sk_path_t encapsulates compound (multiple contour) geometric
    paths consisting of straight line segments, quadratic curves, and
    cubic curves.
*/
typedef struct sk_path_t sk_path_t;
/**
    A sk_picture_t holds recorded canvas drawing commands to be played
    back at a later time.
*/
typedef struct sk_picture_t sk_picture_t;
/**
    A sk_picture_recorder_t holds a sk_canvas_t that records commands
    to create a sk_picture_t.
*/
typedef struct sk_picture_recorder_t sk_picture_recorder_t;
/**
    A sk_bbh_factory_t generates an sk_bbox_hierarchy as a display optimization
    for culling invisible calls recorded by a sk_picture_recorder. It may
    be passed in to sk_picture_recorder_begin_recording_with_bbh_factory,
    typically as an instance of the subclass sk_rtree_factory_t.
*/
typedef struct sk_bbh_factory_t sk_bbh_factory_t;
/**
    A sk_rtree_factory_t generates a sk_rtree as a display optimization
    for culling invisible calls recorded by a sk_picture_recorder_t. Inherits
    from sk_bbh_factory_t.
*/
typedef struct sk_rtree_factory_t sk_rtree_factory_t;
/**
    A sk_shader_t specifies the source color(s) for what is being drawn. If a
    paint has no shader, then the paint's color is used. If the paint
    has a shader, then the shader's color(s) are use instead, but they
    are modulated by the paint's alpha.
*/
typedef struct sk_shader_t sk_shader_t;
/**
    A sk_surface_t holds the destination for drawing to a canvas. For
    raster drawing, the destination is an array of pixels in memory.
    For GPU drawing, the destination is a texture or a framebuffer.
*/
typedef struct sk_surface_t sk_surface_t;
/**
    The sk_region encapsulates the geometric region used to specify
    clipping areas for drawing.
*/
typedef struct sk_region_t sk_region_t;
typedef struct sk_region_iterator_t sk_region_iterator_t;
typedef struct sk_region_cliperator_t sk_region_cliperator_t;
typedef struct sk_region_spanerator_t sk_region_spanerator_t;

typedef enum sk_blendmode_t_ {
    CLEAR_SK_BLENDMODE,
    SRC_SK_BLENDMODE,
    DST_SK_BLENDMODE,
    SRCOVER_SK_BLENDMODE,
    DSTOVER_SK_BLENDMODE,
    SRCIN_SK_BLENDMODE,
    DSTIN_SK_BLENDMODE,
    SRCOUT_SK_BLENDMODE,
    DSTOUT_SK_BLENDMODE,
    SRCATOP_SK_BLENDMODE,
    DSTATOP_SK_BLENDMODE,
    XOR_SK_BLENDMODE,
    PLUS_SK_BLENDMODE,
    MODULATE_SK_BLENDMODE,
    SCREEN_SK_BLENDMODE,
    OVERLAY_SK_BLENDMODE,
    DARKEN_SK_BLENDMODE,
    LIGHTEN_SK_BLENDMODE,
    COLORDODGE_SK_BLENDMODE,
    COLORBURN_SK_BLENDMODE,
    HARDLIGHT_SK_BLENDMODE,
    SOFTLIGHT_SK_BLENDMODE,
    DIFFERENCE_SK_BLENDMODE,
    EXCLUSION_SK_BLENDMODE,
    MULTIPLY_SK_BLENDMODE,
    HUE_SK_BLENDMODE,
    SATURATION_SK_BLENDMODE,
    COLOR_SK_BLENDMODE,
    LUMINOSITY_SK_BLENDMODE,
} sk_blendmode_t;

//////////////////////////////////////////////////////////////////////////////////////////

typedef struct sk_point3_t_ {
    float   x;
    float   y;
    float   z;
} sk_point3_t;

typedef struct sk_ipoint_t_ {
    int32_t   x;
    int32_t   y;
} sk_ipoint_t;

typedef struct sk_size_t_ {
    float   w;
    float   h;
} sk_size_t;

typedef struct sk_isize_t_ {
    int32_t   w;
    int32_t   h;
} sk_isize_t;

typedef struct sk_fontmetrics_t_ {
    uint32_t fFlags;
    float    fTop;
    float    fAscent;
    float    fDescent;
    float    fBottom;
    float    fLeading;
    float    fAvgCharWidth;
    float    fMaxCharWidth;
    float    fXMin;
    float    fXMax;
    float    fXHeight;
    float    fCapHeight;
    float    fUnderlineThickness;
    float    fUnderlinePosition;
    float    fStrikeoutThickness;
    float    fStrikeoutPosition;
} sk_fontmetrics_t;

// Flags for fFlags member of sk_fontmetrics_t
#define FONTMETRICS_FLAGS_UNDERLINE_THICKNESS_IS_VALID (1U << 0)
#define FONTMETRICS_FLAGS_UNDERLINE_POSITION_IS_VALID (1U << 1)

/**
    A lightweight managed string.
*/
typedef struct sk_string_t sk_string_t;
/**

    A sk_bitmap_t is an abstraction that specifies a raster bitmap.
*/
typedef struct sk_bitmap_t sk_bitmap_t;
typedef struct sk_pixmap_t sk_pixmap_t;
typedef struct sk_colorfilter_t sk_colorfilter_t;
typedef struct sk_imagefilter_t sk_imagefilter_t;

typedef struct sk_blender_t sk_blender_t;

/**
   A sk_typeface_t pecifies the typeface and intrinsic style of a font.
    This is used in the paint, along with optionally algorithmic settings like
    textSize, textSkewX, textScaleX, kFakeBoldText_Mask, to specify
    how text appears when drawn (and measured).

    Typeface objects are immutable, and so they can be shared between threads.
*/
typedef struct sk_typeface_t sk_typeface_t;
typedef uint32_t sk_font_table_tag_t;
typedef struct sk_fontmgr_t sk_fontmgr_t;
typedef struct sk_fontstyle_t sk_fontstyle_t;
typedef struct sk_fontstyleset_t sk_fontstyleset_t;
/**
 *  Abstraction layer directly on top of an image codec.
 */
typedef struct sk_codec_t sk_codec_t;
typedef struct sk_colorspace_t sk_colorspace_t;
/**
   Various stream types
*/
typedef struct sk_stream_t sk_stream_t;
typedef struct sk_stream_filestream_t sk_stream_filestream_t;
typedef struct sk_stream_asset_t sk_stream_asset_t;
typedef struct sk_stream_memorystream_t sk_stream_memorystream_t;
typedef struct sk_stream_streamrewindable_t sk_stream_streamrewindable_t;
typedef struct sk_wstream_t sk_wstream_t;
typedef struct sk_wstream_filestream_t sk_wstream_filestream_t;
typedef struct sk_wstream_dynamicmemorystream_t sk_wstream_dynamicmemorystream_t;
/**
   High-level API for creating a document-based canvas.
*/
typedef struct sk_document_t sk_document_t;

typedef enum sk_point_mode_t_ {
    POINTS_SK_POINT_MODE,
    LINES_SK_POINT_MODE,
    POLYGON_SK_POINT_MODE
} sk_point_mode_t;

typedef enum sk_text_align_t_ {
    LEFT_SK_TEXT_ALIGN,
    CENTER_SK_TEXT_ALIGN,
    RIGHT_SK_TEXT_ALIGN
} sk_text_align_t;

typedef enum sk_text_encoding_t_ {
    UTF8_SK_TEXT_ENCODING,
    UTF16_SK_TEXT_ENCODING,
    UTF32_SK_TEXT_ENCODING,
    GLYPH_ID_SK_TEXT_ENCODING
} sk_text_encoding_t;

typedef enum sk_path_filltype_t_ {
    WINDING_SK_PATH_FILLTYPE,
    EVENODD_SK_PATH_FILLTYPE,
    INVERSE_WINDING_SK_PATH_FILLTYPE,
    INVERSE_EVENODD_SK_PATH_FILLTYPE,
} sk_path_filltype_t;

typedef enum sk_font_style_slant_t_ {
    UPRIGHT_SK_FONT_STYLE_SLANT = 0,
    ITALIC_SK_FONT_STYLE_SLANT  = 1,
    OBLIQUE_SK_FONT_STYLE_SLANT = 2,
} sk_font_style_slant_t;

typedef enum sk_color_channel_t_ {
    R_SK_COLOR_CHANNEL,
    G_SK_COLOR_CHANNEL,
    B_SK_COLOR_CHANNEL,
    A_SK_COLOR_CHANNEL,
} sk_color_channel_t;

/**
    The logical operations that can be performed when combining two regions.
*/
typedef enum sk_region_op_t_ {
    DIFFERENCE_SK_REGION_OP,          //!< subtract the op region from the first region
    INTERSECT_SK_REGION_OP,           //!< intersect the two regions
    UNION_SK_REGION_OP,               //!< union (inclusive-or) the two regions
    XOR_SK_REGION_OP,                 //!< exclusive-or the two regions
    REVERSE_DIFFERENCE_SK_REGION_OP,  //!< subtract the first region from the op region
    REPLACE_SK_REGION_OP,             //!< replace the dst region with the op region
} sk_region_op_t;

typedef enum sk_clipop_t_ {
    DIFFERENCE_SK_CLIPOP,
    INTERSECT_SK_CLIPOP,
} sk_clipop_t;

/**
 *  Enum describing format of encoded data.
 */
typedef enum sk_encoded_image_format_t_ {
    BMP_SK_ENCODED_FORMAT,
    GIF_SK_ENCODED_FORMAT,
    ICO_SK_ENCODED_FORMAT,
    JPEG_SK_ENCODED_FORMAT,
    PNG_SK_ENCODED_FORMAT,
    WBMP_SK_ENCODED_FORMAT,
    WEBP_SK_ENCODED_FORMAT,
    PKM_SK_ENCODED_FORMAT,
    KTX_SK_ENCODED_FORMAT,
    ASTC_SK_ENCODED_FORMAT,
    DNG_SK_ENCODED_FORMAT,
    HEIF_SK_ENCODED_FORMAT,
    AVIF_SK_ENCODED_FORMAT,
    JPEGXL_SK_ENCODED_FORMAT,
} sk_encoded_image_format_t;

typedef enum sk_encodedorigin_t_ {
    TOP_LEFT_SK_ENCODED_ORIGIN     = 1, // Default
    TOP_RIGHT_SK_ENCODED_ORIGIN    = 2, // Reflected across y-axis
    BOTTOM_RIGHT_SK_ENCODED_ORIGIN = 3, // Rotated 180
    BOTTOM_LEFT_SK_ENCODED_ORIGIN  = 4, // Reflected across x-axis
    LEFT_TOP_SK_ENCODED_ORIGIN     = 5, // Reflected across x-axis, Rotated 90 CCW
    RIGHT_TOP_SK_ENCODED_ORIGIN    = 6, // Rotated 90 CW
    RIGHT_BOTTOM_SK_ENCODED_ORIGIN = 7, // Reflected across x-axis, Rotated 90 CW
    LEFT_BOTTOM_SK_ENCODED_ORIGIN  = 8, // Rotated 90 CCW
    DEFAULT_SK_ENCODED_ORIGIN      = TOP_LEFT_SK_ENCODED_ORIGIN,
} sk_encodedorigin_t;

typedef enum sk_codec_result_t_ {
    SUCCESS_SK_CODEC_RESULT,
    INCOMPLETE_INPUT_SK_CODEC_RESULT,
    ERROR_IN_INPUT_SK_CODEC_RESULT,
    INVALID_CONVERSION_SK_CODEC_RESULT,
    INVALID_SCALE_SK_CODEC_RESULT,
    INVALID_PARAMETERS_SK_CODEC_RESULT,
    INVALID_INPUT_SK_CODEC_RESULT,
    COULD_NOT_REWIND_SK_CODEC_RESULT,
    INTERNAL_ERROR_SK_CODEC_RESULT,
    UNIMPLEMENTED_SK_CODEC_RESULT,
} sk_codec_result_t;

typedef enum sk_codec_zero_initialized_t_ {
    YES_SK_CODEC_ZERO_INITIALIZED,
    NO_SK_CODEC_ZERO_INITIALIZED,
} sk_codec_zero_initialized_t;

typedef struct sk_codec_options_t_ {
    sk_codec_zero_initialized_t fZeroInitialized;
    sk_irect_t* fSubset;
    int fFrameIndex;
    int fPriorFrame;
} sk_codec_options_t;

typedef enum sk_codec_scanline_order_t_ {
    TOP_DOWN_SK_CODEC_SCANLINE_ORDER,
    BOTTOM_UP_SK_CODEC_SCANLINE_ORDER,
} sk_codec_scanline_order_t;

// The verbs that can be foudn on a path
typedef enum sk_path_verb_t_ {
    MOVE_SK_PATH_VERB,
    LINE_SK_PATH_VERB,
    QUAD_SK_PATH_VERB,
    CONIC_SK_PATH_VERB,
    CUBIC_SK_PATH_VERB,
    CLOSE_SK_PATH_VERB,
    DONE_SK_PATH_VERB
} sk_path_verb_t;

typedef struct sk_path_iterator_t sk_path_iterator_t;
typedef struct sk_path_rawiterator_t sk_path_rawiterator_t;

typedef enum sk_path_add_mode_t_ {
    APPEND_SK_PATH_ADD_MODE,
    EXTEND_SK_PATH_ADD_MODE,
} sk_path_add_mode_t;

typedef enum sk_path_segment_mask_t_ {
    LINE_SK_PATH_SEGMENT_MASK  = 1 << 0,
    QUAD_SK_PATH_SEGMENT_MASK  = 1 << 1,
    CONIC_SK_PATH_SEGMENT_MASK = 1 << 2,
    CUBIC_SK_PATH_SEGMENT_MASK = 1 << 3,
} sk_path_segment_mask_t;

typedef enum sk_path_effect_1d_style_t_ {
    TRANSLATE_SK_PATH_EFFECT_1D_STYLE,
    ROTATE_SK_PATH_EFFECT_1D_STYLE,
    MORPH_SK_PATH_EFFECT_1D_STYLE,
} sk_path_effect_1d_style_t;

typedef enum sk_path_effect_trim_mode_t_ {
    NORMAL_SK_PATH_EFFECT_TRIM_MODE,
    INVERTED_SK_PATH_EFFECT_TRIM_MODE,
} sk_path_effect_trim_mode_t;

typedef struct sk_path_effect_t sk_path_effect_t;

typedef enum sk_stroke_cap_t_ {
    BUTT_SK_STROKE_CAP,
    ROUND_SK_STROKE_CAP,
    SQUARE_SK_STROKE_CAP
} sk_stroke_cap_t;

typedef enum sk_stroke_join_t_ {
    MITER_SK_STROKE_JOIN,
    ROUND_SK_STROKE_JOIN,
    BEVEL_SK_STROKE_JOIN
} sk_stroke_join_t;

typedef enum sk_shader_tilemode_t_ {
    CLAMP_SK_SHADER_TILEMODE,
    REPEAT_SK_SHADER_TILEMODE,
    MIRROR_SK_SHADER_TILEMODE,
    DECAL_SK_SHADER_TILEMODE,
} sk_shader_tilemode_t;

typedef enum sk_blurstyle_t_ {
    NORMAL_SK_BLUR_STYLE,   //!< fuzzy inside and outside
    SOLID_SK_BLUR_STYLE,    //!< solid inside, fuzzy outside
    OUTER_SK_BLUR_STYLE,    //!< nothing inside, fuzzy outside
    INNER_SK_BLUR_STYLE,    //!< fuzzy inside, nothing outside
} sk_blurstyle_t;

typedef enum sk_path_direction_t_ {
    CW_SK_PATH_DIRECTION,
    CCW_SK_PATH_DIRECTION,
} sk_path_direction_t;

typedef enum sk_path_arc_size_t_ {
    SMALL_SK_PATH_ARC_SIZE,
    LARGE_SK_PATH_ARC_SIZE,
} sk_path_arc_size_t;

typedef enum sk_paint_style_t_ {
    FILL_SK_PAINT_STYLE,
    STROKE_SK_PAINT_STYLE,
    STROKE_AND_FILL_SK_PAINT_STYLE,
} sk_paint_style_t;

typedef enum sk_font_hinting_t_ {
    NONE_SK_FONT_HINTING,
    SLIGHT_SK_FONT_HINTING,
    NORMAL_SK_FONT_HINTING,
    FULL_SK_FONT_HINTING,
} sk_font_hinting_t;

typedef enum sk_font_edging_t_ {
    ALIAS_SK_FONT_EDGING,
    ANTIALIAS_SK_FONT_EDGING,
    SUBPIXEL_ANTIALIAS_SK_FONT_EDGING,
} sk_font_edging_t;

typedef struct sk_pixelref_factory_t sk_pixelref_factory_t;

typedef enum gr_surfaceorigin_t_ {
    TOP_LEFT_GR_SURFACE_ORIGIN,
    BOTTOM_LEFT_GR_SURFACE_ORIGIN,
} gr_surfaceorigin_t;

typedef struct gr_context_options_t_ {
    bool      fAvoidStencilBuffers;
    int       fRuntimeProgramCacheSize;
    size_t    fGlyphCacheTextureMaximumBytes;
    bool      fAllowPathMaskCaching;
    bool      fDoManualMipmapping;
    int       fBufferMapThreshold;
} gr_context_options_t;

typedef intptr_t gr_backendobject_t;

typedef struct gr_backendrendertarget_t gr_backendrendertarget_t;
typedef struct gr_backendtexture_t gr_backendtexture_t;

typedef struct gr_direct_context_t gr_direct_context_t;
typedef struct gr_recording_context_t gr_recording_context_t;

typedef enum gr_backend_t_ {
    OPENGL_GR_BACKEND,
    VULKAN_GR_BACKEND,
    METAL_GR_BACKEND,
    DIRECT3D_GR_BACKEND,
    DAWN_GR_BACKEND,
} gr_backend_t;

typedef intptr_t gr_backendcontext_t;

typedef struct gr_glinterface_t gr_glinterface_t;

typedef void (*gr_gl_func_ptr)(void);
typedef gr_gl_func_ptr (*gr_gl_get_proc)(void* ctx, const char* name);

typedef struct gr_gl_textureinfo_t_ {
    unsigned int fTarget;
    unsigned int fID;
    unsigned int fFormat;
    bool fProtected;
} gr_gl_textureinfo_t;

typedef struct gr_gl_framebufferinfo_t_ {
    unsigned int fFBOID;
    unsigned int fFormat;
    bool fProtected;
} gr_gl_framebufferinfo_t;

typedef struct vk_instance_t vk_instance_t;
typedef struct gr_vkinterface_t gr_vkinterface_t;
typedef struct vk_physical_device_t vk_physical_device_t;
typedef struct vk_physical_device_features_t vk_physical_device_features_t;
typedef struct vk_physical_device_features_2_t vk_physical_device_features_2_t;
typedef struct vk_device_t vk_device_t;
typedef struct vk_queue_t vk_queue_t;

typedef struct gr_vk_extensions_t gr_vk_extensions_t;
typedef struct gr_vk_memory_allocator_t gr_vk_memory_allocator_t;

typedef VKAPI_ATTR void (VKAPI_CALL *gr_vk_func_ptr)(void);
typedef gr_vk_func_ptr (*gr_vk_get_proc)(void* ctx, const char* name, vk_instance_t* instance, vk_device_t* device);

typedef struct gr_vk_backendcontext_t_ {
    vk_instance_t*                          fInstance;
    vk_physical_device_t*                   fPhysicalDevice;
    vk_device_t*                            fDevice;
    vk_queue_t*                             fQueue;
    uint32_t                                fGraphicsQueueIndex;
    uint32_t                                fMinAPIVersion;
    uint32_t                                fInstanceVersion;
    uint32_t                                fMaxAPIVersion;
    uint32_t                                fExtensions;
    const gr_vk_extensions_t*               fVkExtensions;
    uint32_t                                fFeatures;
    const vk_physical_device_features_t*    fDeviceFeatures;
    const vk_physical_device_features_2_t*  fDeviceFeatures2;
    gr_vk_memory_allocator_t*               fMemoryAllocator;
    gr_vk_get_proc                          fGetProc;
    void*                                   fGetProcUserData;
    bool                                    fOwnsInstanceAndDevice;
    bool                                    fProtectedContext;
} gr_vk_backendcontext_t;

typedef intptr_t gr_vk_backendmemory_t;

typedef struct gr_vk_alloc_t_ {
    uint64_t               fMemory;
    uint64_t               fOffset;
    uint64_t               fSize;
    uint32_t               fFlags;
    gr_vk_backendmemory_t  fBackendMemory;
    bool                   _private_fUsesSystemHeap;
} gr_vk_alloc_t;

typedef struct gr_vk_ycbcrconversioninfo_t_ {
    uint32_t  fFormat;
    uint64_t  fExternalFormat;
    uint32_t  fYcbcrModel;
    uint32_t  fYcbcrRange;
    uint32_t  fXChromaOffset;
    uint32_t  fYChromaOffset;
    uint32_t  fChromaFilter;
    uint32_t  fForceExplicitReconstruction;
    uint32_t  fFormatFeatures;
} gr_vk_ycbcrconversioninfo_t;

typedef struct gr_vk_imageinfo_t_ {
    uint64_t                        fImage;
    gr_vk_alloc_t                   fAlloc;
    uint32_t                        fImageTiling;
    uint32_t                        fImageLayout;
    uint32_t                        fFormat;
    uint32_t                        fImageUsageFlags;
    uint32_t                        fSampleCount;
    uint32_t                        fLevelCount;
    uint32_t                        fCurrentQueueFamily;
    bool                            fProtected;
    gr_vk_ycbcrconversioninfo_t     fYcbcrConversionInfo;
    uint32_t                        fSharingMode;
} gr_vk_imageinfo_t;

typedef struct vk_instance_t vk_instance_t;
typedef struct vk_physical_device_t vk_physical_device_t;
typedef struct vk_device_t vk_device_t;
typedef struct vk_queue_t vk_queue_t;

#define gr_mtl_handle_t const void*

typedef struct gr_mtl_textureinfo_t_ {
    const void* fTexture;
} gr_mtl_textureinfo_t;

typedef enum sk_pathop_t_ {
    DIFFERENCE_SK_PATHOP,
    INTERSECT_SK_PATHOP,
    UNION_SK_PATHOP,
    XOR_SK_PATHOP,
    REVERSE_DIFFERENCE_SK_PATHOP,
} sk_pathop_t;

typedef struct sk_opbuilder_t sk_opbuilder_t;

typedef enum sk_lattice_recttype_t_ {
    DEFAULT_SK_LATTICE_RECT_TYPE,
    TRANSPARENT_SK_LATTICE_RECT_TYPE,
    FIXED_COLOR_SK_LATTICE_RECT_TYPE,
} sk_lattice_recttype_t;

typedef struct sk_lattice_t_ {
    const int* fXDivs;
    const int* fYDivs;
    const sk_lattice_recttype_t* fRectTypes;
    int fXCount;
    int fYCount;
    const sk_irect_t* fBounds;
    const sk_color_t* fColors;
} sk_lattice_t;

typedef struct sk_pathmeasure_t sk_pathmeasure_t;

typedef enum sk_pathmeasure_matrixflags_t_ {
    GET_POSITION_SK_PATHMEASURE_MATRIXFLAGS = 0x01,
    GET_TANGENT_SK_PATHMEASURE_MATRIXFLAGS = 0x02,
    GET_POS_AND_TAN_SK_PATHMEASURE_MATRIXFLAGS = GET_POSITION_SK_PATHMEASURE_MATRIXFLAGS | GET_TANGENT_SK_PATHMEASURE_MATRIXFLAGS,
} sk_pathmeasure_matrixflags_t;

typedef void (*sk_bitmap_release_proc)(void* addr, void* context);

typedef void (*sk_data_release_proc)(const void* ptr, void* context);

typedef void (*sk_image_raster_release_proc)(const void* addr, void* context);
typedef void (*sk_image_texture_release_proc)(void* context);

typedef void (*sk_surface_raster_release_proc)(void* addr, void* context);

typedef void (*sk_glyph_path_proc)(const sk_path_t* pathOrNull, const sk_matrix_t* matrix, void* context);

typedef enum sk_image_caching_hint_t_ {
    ALLOW_SK_IMAGE_CACHING_HINT,
    DISALLOW_SK_IMAGE_CACHING_HINT,
} sk_image_caching_hint_t;

typedef enum sk_bitmap_allocflags_t_ {
    NONE_SK_BITMAP_ALLOC_FLAGS = 0,
    ZERO_PIXELS_SK_BITMAP_ALLOC_FLAGS = 1 << 0,
} sk_bitmap_allocflags_t;

typedef struct sk_time_datetime_t_ {
    int16_t  fTimeZoneMinutes;
    uint16_t fYear;
    uint8_t  fMonth;
    uint8_t  fDayOfWeek;
    uint8_t  fDay;
    uint8_t  fHour;
    uint8_t  fMinute;
    uint8_t  fSecond;
} sk_time_datetime_t;

typedef struct sk_document_pdf_metadata_t_ {
    sk_string_t*        fTitle;
    sk_string_t*        fAuthor;
    sk_string_t*        fSubject;
    sk_string_t*        fKeywords;
    sk_string_t*        fCreator;
    sk_string_t*        fProducer;
    sk_time_datetime_t* fCreation;
    sk_time_datetime_t* fModified;
    float               fRasterDPI;
    bool                fPDFA;
    int                 fEncodingQuality;
} sk_document_pdf_metadata_t;

typedef struct sk_imageinfo_t_ {
    sk_colorspace_t* colorspace;
    int32_t          width;
    int32_t          height;
    sk_colortype_t   colorType;
    sk_alphatype_t   alphaType;
} sk_imageinfo_t;

typedef enum sk_codecanimation_disposalmethod_t_ {
    KEEP_SK_CODEC_ANIMATION_DISPOSAL_METHOD               = 1,
    RESTORE_BG_COLOR_SK_CODEC_ANIMATION_DISPOSAL_METHOD   = 2,
    RESTORE_PREVIOUS_SK_CODEC_ANIMATION_DISPOSAL_METHOD   = 3,
} sk_codecanimation_disposalmethod_t;

typedef enum sk_codecanimation_blend_t_ {
    SRC_OVER_SK_CODEC_ANIMATION_BLEND   = 0,
    SRC_SK_CODEC_ANIMATION_BLEND        = 1,
} sk_codecanimation_blend_t;

typedef struct sk_codec_frameinfo_t_ {
    int fRequiredFrame;
    int fDuration;
    bool fFullyReceived;
    sk_alphatype_t fAlphaType;
    bool fHasAlphaWithinBounds;
    sk_codecanimation_disposalmethod_t fDisposalMethod;
    sk_codecanimation_blend_t fBlend;
    sk_irect_t fFrameRect;
} sk_codec_frameinfo_t;

typedef struct sk_svgcanvas_t sk_svgcanvas_t;

typedef enum sk_vertices_vertex_mode_t_ {
    TRIANGLES_SK_VERTICES_VERTEX_MODE,
    TRIANGLE_STRIP_SK_VERTICES_VERTEX_MODE,
    TRIANGLE_FAN_SK_VERTICES_VERTEX_MODE,
} sk_vertices_vertex_mode_t;

typedef struct sk_vertices_t sk_vertices_t;

typedef struct sk_colorspace_transfer_fn_t_ {
    float fG;
    float fA;
    float fB;
    float fC;
    float fD;
    float fE;
    float fF;
} sk_colorspace_transfer_fn_t;

typedef struct sk_colorspace_primaries_t_ {
    float fRX;
    float fRY;
    float fGX;
    float fGY;
    float fBX;
    float fBY;
    float fWX;
    float fWY;
} sk_colorspace_primaries_t;

typedef struct sk_colorspace_xyz_t_ {
    float fM00;
    float fM01;
    float fM02;
    float fM10;
    float fM11;
    float fM12;
    float fM20;
    float fM21;
    float fM22;
} sk_colorspace_xyz_t;

typedef struct sk_colorspace_icc_profile_t sk_colorspace_icc_profile_t;

typedef enum sk_highcontrastconfig_invertstyle_t_ {
    NO_INVERT_SK_HIGH_CONTRAST_CONFIG_INVERT_STYLE,
    INVERT_BRIGHTNESS_SK_HIGH_CONTRAST_CONFIG_INVERT_STYLE,
    INVERT_LIGHTNESS_SK_HIGH_CONTRAST_CONFIG_INVERT_STYLE,
} sk_highcontrastconfig_invertstyle_t;

typedef struct sk_highcontrastconfig_t_ {
    bool fGrayscale;
    sk_highcontrastconfig_invertstyle_t fInvertStyle;
    float fContrast;
} sk_highcontrastconfig_t;

typedef enum sk_pngencoder_filterflags_t_ {
    ZERO_SK_PNGENCODER_FILTER_FLAGS  = 0x00,
    NONE_SK_PNGENCODER_FILTER_FLAGS  = 0x08,
    SUB_SK_PNGENCODER_FILTER_FLAGS   = 0x10,
    UP_SK_PNGENCODER_FILTER_FLAGS    = 0x20,
    AVG_SK_PNGENCODER_FILTER_FLAGS   = 0x40,
    PAETH_SK_PNGENCODER_FILTER_FLAGS = 0x80,
    ALL_SK_PNGENCODER_FILTER_FLAGS   = NONE_SK_PNGENCODER_FILTER_FLAGS |
                                       SUB_SK_PNGENCODER_FILTER_FLAGS |
                                       UP_SK_PNGENCODER_FILTER_FLAGS |
                                       AVG_SK_PNGENCODER_FILTER_FLAGS |
                                       PAETH_SK_PNGENCODER_FILTER_FLAGS,
} sk_pngencoder_filterflags_t;

typedef struct sk_pngencoder_options_t_ {
    sk_pngencoder_filterflags_t fFilterFlags;
    int fZLibLevel;
    void* fComments;
    const sk_colorspace_icc_profile_t* fICCProfile;
    const char* fICCProfileDescription;
} sk_pngencoder_options_t;

typedef enum sk_jpegencoder_downsample_t_ {
    DOWNSAMPLE_420_SK_JPEGENCODER_DOWNSAMPLE,
    DOWNSAMPLE_422_SK_JPEGENCODER_DOWNSAMPLE,
    DOWNSAMPLE_444_SK_JPEGENCODER_DOWNSAMPLE,
} sk_jpegencoder_downsample_t;

typedef enum sk_jpegencoder_alphaoption_t_ {
    IGNORE_SK_JPEGENCODER_ALPHA_OPTION,
    BLEND_ON_BLACK_SK_JPEGENCODER_ALPHA_OPTION,
} sk_jpegencoder_alphaoption_t;

typedef struct sk_jpegencoder_options_t_ {
    int fQuality;
    sk_jpegencoder_downsample_t fDownsample;
    sk_jpegencoder_alphaoption_t fAlphaOption;
    const sk_data_t* xmpMetadata;
    const sk_colorspace_icc_profile_t* fICCProfile;
    const char* fICCProfileDescription;
} sk_jpegencoder_options_t;

typedef enum sk_webpencoder_compression_t_ {
    LOSSY_SK_WEBPENCODER_COMPTRESSION,
    LOSSLESS_SK_WEBPENCODER_COMPTRESSION,
} sk_webpencoder_compression_t;

typedef struct sk_webpencoder_options_t_ {
    sk_webpencoder_compression_t fCompression;
    float fQuality;
    const sk_colorspace_icc_profile_t* fICCProfile;
    const char* fICCProfileDescription;
} sk_webpencoder_options_t;

typedef struct sk_rrect_t sk_rrect_t;

typedef enum sk_rrect_type_t_ {
    EMPTY_SK_RRECT_TYPE,
    RECT_SK_RRECT_TYPE,
    OVAL_SK_RRECT_TYPE,
    SIMPLE_SK_RRECT_TYPE,
    NINE_PATCH_SK_RRECT_TYPE,
    COMPLEX_SK_RRECT_TYPE,
} sk_rrect_type_t;

typedef enum sk_rrect_corner_t_ {
    UPPER_LEFT_SK_RRECT_CORNER,
    UPPER_RIGHT_SK_RRECT_CORNER,
    LOWER_RIGHT_SK_RRECT_CORNER,
    LOWER_LEFT_SK_RRECT_CORNER,
} sk_rrect_corner_t;

typedef struct sk_textblob_t sk_textblob_t;
typedef struct sk_textblob_builder_t sk_textblob_builder_t;

typedef struct sk_textblob_builder_runbuffer_t_ {
    void* glyphs;
    void* pos;
    void* utf8text;
    void* clusters;
} sk_textblob_builder_runbuffer_t;

typedef struct sk_rsxform_t_ {
    float fSCos;
    float fSSin;
    float fTX;
    float fTY;
} sk_rsxform_t;

typedef struct sk_tracememorydump_t sk_tracememorydump_t;

typedef struct sk_runtimeeffect_t sk_runtimeeffect_t;

typedef enum sk_runtimeeffect_uniform_type_t_ {
    FLOAT_SK_RUNTIMEEFFECT_UNIFORM_TYPE,
    FLOAT2_SK_RUNTIMEEFFECT_UNIFORM_TYPE,
    FLOAT3_SK_RUNTIMEEFFECT_UNIFORM_TYPE,
    FLOAT4_SK_RUNTIMEEFFECT_UNIFORM_TYPE,
    FLOAT2X2_SK_RUNTIMEEFFECT_UNIFORM_TYPE,
    FLOAT3X3_SK_RUNTIMEEFFECT_UNIFORM_TYPE,
    FLOAT4X4_SK_RUNTIMEEFFECT_UNIFORM_TYPE,
    INT_SK_RUNTIMEEFFECT_UNIFORM_TYPE,
    INT2_SK_RUNTIMEEFFECT_UNIFORM_TYPE,
    INT3_SK_RUNTIMEEFFECT_UNIFORM_TYPE,
    INT4_SK_RUNTIMEEFFECT_UNIFORM_TYPE,
} sk_runtimeeffect_uniform_type_t;

typedef enum sk_runtimeeffect_child_type_t_ {
    SHADER_SK_RUNTIMEEFFECT_CHILD_TYPE,
    COLOR_FILTER_SK_RUNTIMEEFFECT_CHILD_TYPE,
    BLENDER_SK_RUNTIMEEFFECT_CHILD_TYPE,
} sk_runtimeeffect_child_type_t;

typedef enum sk_runtimeeffect_uniform_flags_t_ {
    NONE_SK_RUNTIMEEFFECT_UNIFORM_FLAGS           = 0x00,
    ARRAY_SK_RUNTIMEEFFECT_UNIFORM_FLAGS          = 0x01,
    COLOR_SK_RUNTIMEEFFECT_UNIFORM_FLAGS          = 0x02,
    VERTEX_SK_RUNTIMEEFFECT_UNIFORM_FLAGS         = 0x04,
    FRAGMENT_SK_RUNTIMEEFFECT_UNIFORM_FLAGS       = 0x08,
    HALF_PRECISION_SK_RUNTIMEEFFECT_UNIFORM_FLAGS = 0x10,
} sk_runtimeeffect_uniform_flags_t;

typedef struct sk_runtimeeffect_uniform_t_ {
    const char* fName;
    size_t fNameLength;
    size_t fOffset;
    sk_runtimeeffect_uniform_type_t fType;
    int fCount;
    sk_runtimeeffect_uniform_flags_t fFlags;
} sk_runtimeeffect_uniform_t;

typedef struct sk_runtimeeffect_child_t_ {
    const char* fName;
    size_t fNameLength;
    sk_runtimeeffect_child_type_t fType;
    int fIndex;
} sk_runtimeeffect_child_t;

typedef enum sk_filter_mode_t_ {
    NEAREST_SK_FILTER_MODE,
    LINEAR_SK_FILTER_MODE,
} sk_filter_mode_t;

typedef enum sk_mipmap_mode_t_ {
    NONE_SK_MIPMAP_MODE,
    NEAREST_SK_MIPMAP_MODE,
    LINEAR_SK_MIPMAP_MODE,
} sk_mipmap_mode_t;

typedef struct sk_cubic_resampler_t_ {
    float fB;
    float fC;
} sk_cubic_resampler_t;

typedef struct sk_sampling_options_t_ {
    int fMaxAniso;
    bool fUseCubic;
    sk_cubic_resampler_t fCubic;
    sk_filter_mode_t fFilter;
    sk_mipmap_mode_t fMipmap;
} sk_sampling_options_t;

/*
 * Skottie Animation
 */
typedef struct skottie_animation_t skottie_animation_t;
typedef struct skottie_animation_builder_t skottie_animation_builder_t;
typedef struct skottie_resource_provider_t skottie_resource_provider_t;
typedef struct skottie_property_observer_t skottie_property_observer_t;
typedef struct skottie_logger_t skottie_logger_t;
typedef struct skottie_marker_observer_t skottie_marker_observer_t;

typedef struct sksg_invalidation_controller_t sksg_invalidation_controller_t;

typedef enum skottie_animation_renderflags_t_ {
    SKIP_TOP_LEVEL_ISOLATION = 0x01,
    DISABLE_TOP_LEVEL_CLIPPING = 0x02,
} skottie_animation_renderflags_t;

typedef enum skottie_animation_builder_flags_t_ {
    NONE_SKOTTIE_ANIMATION_BUILDER_FLAGS = 0,
    DEFER_IMAGE_LOADING_SKOTTIE_ANIMATION_BUILDER_FLAGS = 0x01,
    PREFER_EMBEDDED_FONTS_SKOTTIE_ANIMATION_BUILDER_FLAGS = 0x02,
} skottie_animation_builder_flags_t;

typedef struct skottie_animation_builder_stats_t_ {
    float fTotalLoadTimeMS;
    float fJsonParseTimeMS;
    float fSceneParseTimeMS;
    size_t fJsonSize;
    size_t fAnimatorCount;
} skottie_animation_builder_stats_t;

typedef struct skresources_image_asset_t skresources_image_asset_t;
typedef struct skresources_multi_frame_image_asset_t skresources_multi_frame_image_asset_t;
typedef struct skresources_external_track_asset_t skresources_external_track_asset_t;

typedef struct skresources_resource_provider_t skresources_resource_provider_t;

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\gr_context.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef gr_context_DEFINED
#define gr_context_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

// GrRecordingContext

SK_C_API void gr_recording_context_unref(gr_recording_context_t* context);
SK_C_API int gr_recording_context_get_max_surface_sample_count_for_color_type(gr_recording_context_t* context, sk_colortype_t colorType);
SK_C_API gr_backend_t gr_recording_context_get_backend(gr_recording_context_t* context);
SK_C_API bool gr_recording_context_is_abandoned(gr_recording_context_t* context);
SK_C_API int gr_recording_context_max_texture_size(gr_recording_context_t* context);
SK_C_API int gr_recording_context_max_render_target_size(gr_recording_context_t* context);

// GrDirectContext

SK_C_API gr_direct_context_t* gr_direct_context_make_gl(const gr_glinterface_t* glInterface);
SK_C_API gr_direct_context_t* gr_direct_context_make_gl_with_options(const gr_glinterface_t* glInterface, const gr_context_options_t* options);
SK_C_API gr_direct_context_t* gr_direct_context_make_vulkan(const gr_vk_backendcontext_t vkBackendContext);
SK_C_API gr_direct_context_t* gr_direct_context_make_vulkan_with_options(const gr_vk_backendcontext_t vkBackendContext, const gr_context_options_t* options);
SK_C_API gr_direct_context_t* gr_direct_context_make_metal(void* device, void* queue);
SK_C_API gr_direct_context_t* gr_direct_context_make_metal_with_options(void* device, void* queue, const gr_context_options_t* options);

// TODO: the overloads with GrContextOptions

SK_C_API bool gr_direct_context_is_abandoned(gr_direct_context_t* context);
SK_C_API void gr_direct_context_abandon_context(gr_direct_context_t* context);
SK_C_API void gr_direct_context_release_resources_and_abandon_context(gr_direct_context_t* context);
SK_C_API size_t gr_direct_context_get_resource_cache_limit(gr_direct_context_t* context);
SK_C_API void gr_direct_context_set_resource_cache_limit(gr_direct_context_t* context, size_t maxResourceBytes);
SK_C_API void gr_direct_context_get_resource_cache_usage(gr_direct_context_t* context, int* maxResources, size_t* maxResourceBytes);
SK_C_API void gr_direct_context_flush(gr_direct_context_t* context);
SK_C_API bool gr_direct_context_submit(gr_direct_context_t* context, bool syncCpu);
SK_C_API void gr_direct_context_flush_and_submit(gr_direct_context_t* context, bool syncCpu);
SK_C_API void gr_direct_context_flush_image(gr_direct_context_t* context, const sk_image_t* image);
SK_C_API void gr_direct_context_flush_surface(gr_direct_context_t* context, sk_surface_t* surface);
SK_C_API void gr_direct_context_reset_context(gr_direct_context_t* context, uint32_t state);
SK_C_API void gr_direct_context_dump_memory_statistics(const gr_direct_context_t* context, sk_tracememorydump_t* dump);
SK_C_API void gr_direct_context_free_gpu_resources(gr_direct_context_t* context);
SK_C_API void gr_direct_context_perform_deferred_cleanup(gr_direct_context_t* context, long long ms);
SK_C_API void gr_direct_context_purge_unlocked_resources_bytes(gr_direct_context_t* context, size_t bytesToPurge, bool preferScratchResources);
SK_C_API void gr_direct_context_purge_unlocked_resources(gr_direct_context_t* context, bool scratchResourcesOnly);


// GrGLInterface

SK_C_API const gr_glinterface_t* gr_glinterface_create_native_interface(void);
SK_C_API const gr_glinterface_t* gr_glinterface_assemble_interface(void* ctx, gr_gl_get_proc get);
SK_C_API const gr_glinterface_t* gr_glinterface_assemble_gl_interface(void* ctx, gr_gl_get_proc get);
SK_C_API const gr_glinterface_t* gr_glinterface_assemble_gles_interface(void* ctx, gr_gl_get_proc get);
SK_C_API const gr_glinterface_t* gr_glinterface_assemble_webgl_interface(void* ctx, gr_gl_get_proc get);

SK_C_API void gr_glinterface_unref(const gr_glinterface_t* glInterface);
SK_C_API bool gr_glinterface_validate(const gr_glinterface_t* glInterface);
SK_C_API bool gr_glinterface_has_extension(const gr_glinterface_t* glInterface, const char* extension);

// GrVkExtensions

SK_C_API gr_vk_extensions_t* gr_vk_extensions_new(void);
SK_C_API void gr_vk_extensions_delete(gr_vk_extensions_t* extensions);
SK_C_API void gr_vk_extensions_init(gr_vk_extensions_t* extensions, gr_vk_get_proc getProc, void* userData, vk_instance_t* instance, vk_physical_device_t* physDev, uint32_t instanceExtensionCount, const char** instanceExtensions, uint32_t deviceExtensionCount, const char** deviceExtensions);
SK_C_API bool gr_vk_extensions_has_extension(gr_vk_extensions_t* extensions, const char* ext, uint32_t minVersion);

// GrBackendTexture

SK_C_API gr_backendtexture_t* gr_backendtexture_new_gl(int width, int height, bool mipmapped, const gr_gl_textureinfo_t* glInfo);
SK_C_API gr_backendtexture_t* gr_backendtexture_new_vulkan(int width, int height, const gr_vk_imageinfo_t* vkInfo);
SK_C_API gr_backendtexture_t* gr_backendtexture_new_metal(int width, int height, bool mipmapped, const gr_mtl_textureinfo_t* mtlInfo);
SK_C_API void gr_backendtexture_delete(gr_backendtexture_t* texture);

SK_C_API bool gr_backendtexture_is_valid(const gr_backendtexture_t* texture);
SK_C_API int gr_backendtexture_get_width(const gr_backendtexture_t* texture);
SK_C_API int gr_backendtexture_get_height(const gr_backendtexture_t* texture);
SK_C_API bool gr_backendtexture_has_mipmaps(const gr_backendtexture_t* texture);
SK_C_API gr_backend_t gr_backendtexture_get_backend(const gr_backendtexture_t* texture);
SK_C_API bool gr_backendtexture_get_gl_textureinfo(const gr_backendtexture_t* texture, gr_gl_textureinfo_t* glInfo);


// GrBackendRenderTarget

SK_C_API gr_backendrendertarget_t* gr_backendrendertarget_new_gl(int width, int height, int samples, int stencils, const gr_gl_framebufferinfo_t* glInfo);
SK_C_API gr_backendrendertarget_t* gr_backendrendertarget_new_vulkan(int width, int height, int samples, const gr_vk_imageinfo_t* vkImageInfo);
SK_C_API gr_backendrendertarget_t* gr_backendrendertarget_new_metal(int width, int height, int samples, const gr_mtl_textureinfo_t* mtlInfo);

SK_C_API void gr_backendrendertarget_delete(gr_backendrendertarget_t* rendertarget);

SK_C_API bool gr_backendrendertarget_is_valid(const gr_backendrendertarget_t* rendertarget);
SK_C_API int gr_backendrendertarget_get_width(const gr_backendrendertarget_t* rendertarget);
SK_C_API int gr_backendrendertarget_get_height(const gr_backendrendertarget_t* rendertarget);
SK_C_API int gr_backendrendertarget_get_samples(const gr_backendrendertarget_t* rendertarget);
SK_C_API int gr_backendrendertarget_get_stencils(const gr_backendrendertarget_t* rendertarget);
SK_C_API gr_backend_t gr_backendrendertarget_get_backend(const gr_backendrendertarget_t* rendertarget);
SK_C_API bool gr_backendrendertarget_get_gl_framebufferinfo(const gr_backendrendertarget_t* rendertarget, gr_gl_framebufferinfo_t* glInfo);


SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_bitmap.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_bitmap_DEFINED
#define sk_bitmap_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void sk_bitmap_destructor(sk_bitmap_t* cbitmap);
SK_C_API sk_bitmap_t* sk_bitmap_new(void);
SK_C_API void sk_bitmap_get_info(sk_bitmap_t* cbitmap, sk_imageinfo_t* info);
SK_C_API void* sk_bitmap_get_pixels(sk_bitmap_t* cbitmap, size_t* length);
SK_C_API size_t sk_bitmap_get_row_bytes(sk_bitmap_t* cbitmap);
SK_C_API size_t sk_bitmap_get_byte_count(sk_bitmap_t* cbitmap);
SK_C_API void sk_bitmap_reset(sk_bitmap_t* cbitmap);
SK_C_API bool sk_bitmap_is_null(sk_bitmap_t* cbitmap);
SK_C_API bool sk_bitmap_is_immutable(sk_bitmap_t* cbitmap);
SK_C_API void sk_bitmap_set_immutable(sk_bitmap_t* cbitmap);
SK_C_API void sk_bitmap_erase(sk_bitmap_t* cbitmap, sk_color_t color);
SK_C_API void sk_bitmap_erase_rect(sk_bitmap_t* cbitmap, sk_color_t color, sk_irect_t* rect);
SK_C_API uint8_t* sk_bitmap_get_addr_8(sk_bitmap_t* cbitmap, int x, int y);
SK_C_API uint16_t* sk_bitmap_get_addr_16(sk_bitmap_t* cbitmap, int x, int y);
SK_C_API uint32_t* sk_bitmap_get_addr_32(sk_bitmap_t* cbitmap, int x, int y);
SK_C_API void* sk_bitmap_get_addr(sk_bitmap_t* cbitmap, int x, int y);
SK_C_API sk_color_t sk_bitmap_get_pixel_color(sk_bitmap_t* cbitmap, int x, int y);
SK_C_API bool sk_bitmap_ready_to_draw(sk_bitmap_t* cbitmap);
SK_C_API void sk_bitmap_get_pixel_colors(sk_bitmap_t* cbitmap, sk_color_t* colors);
SK_C_API bool sk_bitmap_install_pixels(sk_bitmap_t* cbitmap, const sk_imageinfo_t* cinfo, void* pixels, size_t rowBytes, const sk_bitmap_release_proc releaseProc, void* context);
SK_C_API bool sk_bitmap_install_pixels_with_pixmap(sk_bitmap_t* cbitmap, const sk_pixmap_t* cpixmap);
SK_C_API bool sk_bitmap_try_alloc_pixels(sk_bitmap_t* cbitmap, const sk_imageinfo_t* requestedInfo, size_t rowBytes);
SK_C_API bool sk_bitmap_try_alloc_pixels_with_flags(sk_bitmap_t* cbitmap, const sk_imageinfo_t* requestedInfo, uint32_t flags);
SK_C_API void sk_bitmap_set_pixels(sk_bitmap_t* cbitmap, void* pixels);
SK_C_API bool sk_bitmap_peek_pixels(sk_bitmap_t* cbitmap, sk_pixmap_t* cpixmap);
SK_C_API bool sk_bitmap_extract_subset(sk_bitmap_t* cbitmap, sk_bitmap_t* dst, sk_irect_t* subset);
SK_C_API bool sk_bitmap_extract_alpha(sk_bitmap_t* cbitmap, sk_bitmap_t* dst, const sk_paint_t* paint, sk_ipoint_t* offset);
SK_C_API void sk_bitmap_notify_pixels_changed(sk_bitmap_t* cbitmap);
SK_C_API void sk_bitmap_swap(sk_bitmap_t* cbitmap, sk_bitmap_t* cother);
SK_C_API sk_shader_t* sk_bitmap_make_shader(sk_bitmap_t* cbitmap, sk_shader_tilemode_t tmx, sk_shader_tilemode_t tmy, sk_sampling_options_t* sampling, const sk_matrix_t* cmatrix);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_blender.h
/*
 * Copyright 2024 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_blender_DEFINED
#define sk_blender_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void sk_blender_ref(sk_blender_t* blender);
SK_C_API void sk_blender_unref(sk_blender_t* blender);
SK_C_API sk_blender_t* sk_blender_new_mode(sk_blendmode_t mode);
SK_C_API sk_blender_t* sk_blender_new_arithmetic(float k1, float k2, float k3, float k4, bool enforcePremul);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_canvas.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_canvas_DEFINED
#define sk_canvas_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void sk_canvas_destroy(sk_canvas_t* ccanvas);
SK_C_API void sk_canvas_clear(sk_canvas_t* ccanvas, sk_color_t color);
SK_C_API void sk_canvas_clear_color4f(sk_canvas_t* ccanvas, sk_color4f_t color);
SK_C_API void sk_canvas_discard(sk_canvas_t* ccanvas);
SK_C_API int sk_canvas_get_save_count(sk_canvas_t* ccanvas);
SK_C_API void sk_canvas_restore_to_count(sk_canvas_t* ccanvas, int saveCount);
SK_C_API void sk_canvas_draw_color(sk_canvas_t* ccanvas, sk_color_t color, sk_blendmode_t cmode);
SK_C_API void sk_canvas_draw_color4f(sk_canvas_t* ccanvas, sk_color4f_t color, sk_blendmode_t cmode);
SK_C_API void sk_canvas_draw_points(sk_canvas_t* ccanvas, sk_point_mode_t pointMode, size_t count, const sk_point_t points [], const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_point(sk_canvas_t* ccanvas, float x, float y, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_line(sk_canvas_t* ccanvas, float x0, float y0, float x1, float y1, sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_simple_text(sk_canvas_t* ccanvas, const void* text, size_t byte_length, sk_text_encoding_t encoding, float x, float y, const sk_font_t* cfont, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_text_blob (sk_canvas_t* ccanvas, sk_textblob_t* text, float x, float y, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_reset_matrix(sk_canvas_t* ccanvas);
SK_C_API void sk_canvas_set_matrix(sk_canvas_t* ccanvas, const sk_matrix44_t* cmatrix);
SK_C_API void sk_canvas_get_matrix(sk_canvas_t* ccanvas, sk_matrix44_t* cmatrix);
SK_C_API void sk_canvas_draw_round_rect(sk_canvas_t* ccanvas, const sk_rect_t* crect, float rx, float ry, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_clip_rect_with_operation(sk_canvas_t* ccanvas, const sk_rect_t* crect, sk_clipop_t op, bool doAA);
SK_C_API void sk_canvas_clip_path_with_operation(sk_canvas_t* ccanvas, const sk_path_t* cpath, sk_clipop_t op, bool doAA);
SK_C_API void sk_canvas_clip_rrect_with_operation(sk_canvas_t* ccanvas, const sk_rrect_t* crect, sk_clipop_t op, bool doAA);
SK_C_API bool sk_canvas_get_local_clip_bounds(sk_canvas_t* ccanvas, sk_rect_t* cbounds);
SK_C_API bool sk_canvas_get_device_clip_bounds(sk_canvas_t* ccanvas, sk_irect_t* cbounds);
SK_C_API int sk_canvas_save(sk_canvas_t* ccanvas);
SK_C_API int sk_canvas_save_layer(sk_canvas_t* ccanvas, const sk_rect_t* crect, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_restore(sk_canvas_t* ccanvas);
SK_C_API void sk_canvas_translate(sk_canvas_t* ccanvas, float dx, float dy);
SK_C_API void sk_canvas_scale(sk_canvas_t* ccanvas, float sx, float sy);
SK_C_API void sk_canvas_rotate_degrees(sk_canvas_t* ccanvas, float degrees);
SK_C_API void sk_canvas_rotate_radians(sk_canvas_t* ccanvas, float radians);
SK_C_API void sk_canvas_skew(sk_canvas_t* ccanvas, float sx, float sy);
SK_C_API void sk_canvas_concat(sk_canvas_t* ccanvas, const sk_matrix44_t* cmatrix);
SK_C_API bool sk_canvas_quick_reject(sk_canvas_t* ccanvas, const sk_rect_t* crect);
SK_C_API void sk_canvas_clip_region(sk_canvas_t* ccanvas, const sk_region_t* region, sk_clipop_t op);
SK_C_API void sk_canvas_draw_paint(sk_canvas_t* ccanvas, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_region(sk_canvas_t* ccanvas, const sk_region_t* cregion, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_rect(sk_canvas_t* ccanvas, const sk_rect_t* crect, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_rrect(sk_canvas_t* ccanvas, const sk_rrect_t* crect, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_circle(sk_canvas_t* ccanvas, float cx, float cy, float rad, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_oval(sk_canvas_t* ccanvas, const sk_rect_t* crect, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_path(sk_canvas_t* ccanvas, const sk_path_t* cpath, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_image(sk_canvas_t* ccanvas, const sk_image_t* cimage, float x, float y, const sk_sampling_options_t* sampling, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_image_rect(sk_canvas_t* ccanvas, const sk_image_t* cimage, const sk_rect_t* csrcR, const sk_rect_t* cdstR, const sk_sampling_options_t* sampling, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_picture(sk_canvas_t* ccanvas, const sk_picture_t* cpicture, const sk_matrix_t* cmatrix, const sk_paint_t* cpaint);
SK_C_API void sk_canvas_draw_drawable(sk_canvas_t* ccanvas, sk_drawable_t* cdrawable, const sk_matrix_t* cmatrix);
SK_C_API void sk_canvas_flush(sk_canvas_t* ccanvas);
SK_C_API sk_canvas_t* sk_canvas_new_from_bitmap(const sk_bitmap_t* bitmap);
SK_C_API sk_canvas_t* sk_canvas_new_from_raster(const sk_imageinfo_t* cinfo, void* pixels, size_t rowBytes, const sk_surfaceprops_t* props);
SK_C_API void sk_canvas_draw_annotation(sk_canvas_t* t, const sk_rect_t* rect, const char* key, sk_data_t* value);
SK_C_API void sk_canvas_draw_url_annotation(sk_canvas_t* t, const sk_rect_t* rect, sk_data_t* value);
SK_C_API void sk_canvas_draw_named_destination_annotation(sk_canvas_t* t, const sk_point_t* point, sk_data_t* value);
SK_C_API void sk_canvas_draw_link_destination_annotation(sk_canvas_t* t, const sk_rect_t* rect, sk_data_t* value);
SK_C_API void sk_canvas_draw_image_lattice(sk_canvas_t* ccanvas, const sk_image_t* image, const sk_lattice_t* lattice, const sk_rect_t* dst, sk_filter_mode_t mode, const sk_paint_t* paint);
SK_C_API void sk_canvas_draw_image_nine(sk_canvas_t* ccanvas, const sk_image_t* image, const sk_irect_t* center, const sk_rect_t* dst, sk_filter_mode_t mode, const sk_paint_t* paint);
SK_C_API void sk_canvas_draw_vertices(sk_canvas_t* ccanvas, const sk_vertices_t* vertices, sk_blendmode_t mode, const sk_paint_t* paint);
SK_C_API void sk_canvas_draw_arc(sk_canvas_t* ccanvas, const sk_rect_t* oval, float startAngle, float sweepAngle, bool useCenter, const sk_paint_t* paint);
SK_C_API void sk_canvas_draw_drrect(sk_canvas_t* ccanvas, const sk_rrect_t* outer, const sk_rrect_t* inner, const sk_paint_t* paint);
SK_C_API void sk_canvas_draw_atlas(sk_canvas_t* ccanvas, const sk_image_t* atlas, const sk_rsxform_t* xform, const sk_rect_t* tex, const sk_color_t* colors, int count, sk_blendmode_t mode, const sk_sampling_options_t* sampling, const sk_rect_t* cullRect, const sk_paint_t* paint);
SK_C_API void sk_canvas_draw_patch(sk_canvas_t* ccanvas, const sk_point_t* cubics, const sk_color_t* colors, const sk_point_t* texCoords, sk_blendmode_t mode, const sk_paint_t* paint);
SK_C_API bool sk_canvas_is_clip_empty(sk_canvas_t* ccanvas);
SK_C_API bool sk_canvas_is_clip_rect(sk_canvas_t* ccanvas);
SK_C_API sk_nodraw_canvas_t* sk_nodraw_canvas_new(int width, int height);
SK_C_API void sk_nodraw_canvas_destroy(sk_nodraw_canvas_t* t);
SK_C_API sk_nway_canvas_t* sk_nway_canvas_new(int width, int height);
SK_C_API void sk_nway_canvas_destroy(sk_nway_canvas_t* t);
SK_C_API void sk_nway_canvas_add_canvas(sk_nway_canvas_t* t, sk_canvas_t* canvas);
SK_C_API void sk_nway_canvas_remove_canvas(sk_nway_canvas_t* t, sk_canvas_t* canvas);
SK_C_API void sk_nway_canvas_remove_all(sk_nway_canvas_t* t);
SK_C_API sk_overdraw_canvas_t* sk_overdraw_canvas_new(sk_canvas_t* canvas);
SK_C_API void sk_overdraw_canvas_destroy(sk_overdraw_canvas_t* canvas);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_codec.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_codec_DEFINED
#define sk_codec_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API size_t sk_codec_min_buffered_bytes_needed(void);

SK_C_API sk_codec_t* sk_codec_new_from_stream(sk_stream_t* stream, sk_codec_result_t* result);
SK_C_API sk_codec_t* sk_codec_new_from_data(sk_data_t* data);
SK_C_API void sk_codec_destroy(sk_codec_t* codec);
SK_C_API void sk_codec_get_info(sk_codec_t* codec, sk_imageinfo_t* info);
SK_C_API sk_encodedorigin_t sk_codec_get_origin(sk_codec_t* codec);
SK_C_API void sk_codec_get_scaled_dimensions(sk_codec_t* codec, float desiredScale, sk_isize_t* dimensions);
SK_C_API bool sk_codec_get_valid_subset(sk_codec_t* codec, sk_irect_t* desiredSubset);
SK_C_API sk_encoded_image_format_t sk_codec_get_encoded_format(sk_codec_t* codec);
SK_C_API sk_codec_result_t sk_codec_get_pixels(sk_codec_t* codec, const sk_imageinfo_t* info, void* pixels, size_t rowBytes, const sk_codec_options_t* options);
SK_C_API sk_codec_result_t sk_codec_start_incremental_decode(sk_codec_t* codec, const sk_imageinfo_t* info, void* pixels, size_t rowBytes, const sk_codec_options_t* options);
SK_C_API sk_codec_result_t sk_codec_incremental_decode(sk_codec_t* codec, int* rowsDecoded);
SK_C_API sk_codec_result_t sk_codec_start_scanline_decode(sk_codec_t* codec, const sk_imageinfo_t* info, const sk_codec_options_t* options);
SK_C_API int sk_codec_get_scanlines(sk_codec_t* codec, void* dst, int countLines, size_t rowBytes);
SK_C_API bool sk_codec_skip_scanlines(sk_codec_t* codec, int countLines);
SK_C_API sk_codec_scanline_order_t sk_codec_get_scanline_order(sk_codec_t* codec);
SK_C_API int sk_codec_next_scanline(sk_codec_t* codec);
SK_C_API int sk_codec_output_scanline(sk_codec_t* codec, int inputScanline);
SK_C_API int sk_codec_get_frame_count(sk_codec_t* codec);
SK_C_API void sk_codec_get_frame_info(sk_codec_t* codec, sk_codec_frameinfo_t* frameInfo);
SK_C_API bool sk_codec_get_frame_info_for_index(sk_codec_t* codec, int index, sk_codec_frameinfo_t* frameInfo);
SK_C_API int sk_codec_get_repetition_count(sk_codec_t* codec);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_colorfilter.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_colorfilter_DEFINED
#define sk_colorfilter_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void sk_colorfilter_unref(sk_colorfilter_t* filter);
SK_C_API sk_colorfilter_t* sk_colorfilter_new_mode(sk_color_t c, sk_blendmode_t mode);
SK_C_API sk_colorfilter_t* sk_colorfilter_new_lighting(sk_color_t mul, sk_color_t add);
SK_C_API sk_colorfilter_t* sk_colorfilter_new_compose(sk_colorfilter_t* outer, sk_colorfilter_t* inner);
SK_C_API sk_colorfilter_t* sk_colorfilter_new_color_matrix(const float array[20]);
SK_C_API sk_colorfilter_t* sk_colorfilter_new_luma_color(void);
SK_C_API sk_colorfilter_t* sk_colorfilter_new_high_contrast(const sk_highcontrastconfig_t* config);
SK_C_API sk_colorfilter_t* sk_colorfilter_new_table(const uint8_t table[256]);
SK_C_API sk_colorfilter_t* sk_colorfilter_new_table_argb(const uint8_t tableA[256], const uint8_t tableR[256], const uint8_t tableG[256], const uint8_t tableB[256]);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_colorspace.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_colorspace_DEFINED
#define sk_colorspace_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

// TODO: skcms.h has things that may be useful

// sk_colorspace_t

SK_C_API void sk_colorspace_ref(sk_colorspace_t* colorspace);
SK_C_API void sk_colorspace_unref(sk_colorspace_t* colorspace);
SK_C_API sk_colorspace_t* sk_colorspace_new_srgb(void);
SK_C_API sk_colorspace_t* sk_colorspace_new_srgb_linear(void);
SK_C_API sk_colorspace_t* sk_colorspace_new_rgb(const sk_colorspace_transfer_fn_t* transferFn, const sk_colorspace_xyz_t* toXYZD50);
SK_C_API sk_colorspace_t* sk_colorspace_new_icc(const sk_colorspace_icc_profile_t* profile);
SK_C_API void sk_colorspace_to_profile(const sk_colorspace_t* colorspace, sk_colorspace_icc_profile_t* profile);
SK_C_API bool sk_colorspace_gamma_close_to_srgb(const sk_colorspace_t* colorspace);
SK_C_API bool sk_colorspace_gamma_is_linear(const sk_colorspace_t* colorspace);
SK_C_API bool sk_colorspace_is_numerical_transfer_fn(const sk_colorspace_t* colorspace, sk_colorspace_transfer_fn_t* transferFn);
SK_C_API bool sk_colorspace_to_xyzd50(const sk_colorspace_t* colorspace, sk_colorspace_xyz_t* toXYZD50);
SK_C_API sk_colorspace_t* sk_colorspace_make_linear_gamma(const sk_colorspace_t* colorspace);
SK_C_API sk_colorspace_t* sk_colorspace_make_srgb_gamma(const sk_colorspace_t* colorspace);
SK_C_API bool sk_colorspace_is_srgb(const sk_colorspace_t* colorspace);
SK_C_API bool sk_colorspace_equals(const sk_colorspace_t* src, const sk_colorspace_t* dst);

// sk_colorspace_transfer_fn_t

SK_C_API void sk_colorspace_transfer_fn_named_srgb(sk_colorspace_transfer_fn_t* transferFn);
SK_C_API void sk_colorspace_transfer_fn_named_2dot2(sk_colorspace_transfer_fn_t* transferFn);
SK_C_API void sk_colorspace_transfer_fn_named_linear(sk_colorspace_transfer_fn_t* transferFn);
SK_C_API void sk_colorspace_transfer_fn_named_rec2020(sk_colorspace_transfer_fn_t* transferFn);
SK_C_API void sk_colorspace_transfer_fn_named_pq(sk_colorspace_transfer_fn_t* transferFn);
SK_C_API void sk_colorspace_transfer_fn_named_hlg(sk_colorspace_transfer_fn_t* transferFn);
SK_C_API float sk_colorspace_transfer_fn_eval(const sk_colorspace_transfer_fn_t* transferFn, float x);
SK_C_API bool sk_colorspace_transfer_fn_invert(const sk_colorspace_transfer_fn_t* src, sk_colorspace_transfer_fn_t* dst);

// sk_colorspace_primaries_t

SK_C_API bool sk_colorspace_primaries_to_xyzd50(const sk_colorspace_primaries_t* primaries, sk_colorspace_xyz_t* toXYZD50);

// sk_colorspace_xyz_t

SK_C_API void sk_colorspace_xyz_named_srgb(sk_colorspace_xyz_t* xyz);
SK_C_API void sk_colorspace_xyz_named_adobe_rgb(sk_colorspace_xyz_t* xyz);
SK_C_API void sk_colorspace_xyz_named_display_p3(sk_colorspace_xyz_t* xyz);
SK_C_API void sk_colorspace_xyz_named_rec2020(sk_colorspace_xyz_t* xyz);
SK_C_API void sk_colorspace_xyz_named_xyz(sk_colorspace_xyz_t* xyz);
SK_C_API bool sk_colorspace_xyz_invert(const sk_colorspace_xyz_t* src, sk_colorspace_xyz_t* dst);
SK_C_API void sk_colorspace_xyz_concat(const sk_colorspace_xyz_t* a, const sk_colorspace_xyz_t* b, sk_colorspace_xyz_t* result);

// sk_colorspace_icc_profile_t

SK_C_API void sk_colorspace_icc_profile_delete(sk_colorspace_icc_profile_t* profile);
SK_C_API sk_colorspace_icc_profile_t* sk_colorspace_icc_profile_new(void);
SK_C_API bool sk_colorspace_icc_profile_parse(const void* buffer, size_t length, sk_colorspace_icc_profile_t* profile);
SK_C_API const uint8_t* sk_colorspace_icc_profile_get_buffer(const sk_colorspace_icc_profile_t* profile, uint32_t* size);
SK_C_API bool sk_colorspace_icc_profile_get_to_xyzd50(const sk_colorspace_icc_profile_t* profile, sk_colorspace_xyz_t* toXYZD50);

// sk_color4f_t

SK_C_API sk_color_t sk_color4f_to_color(const sk_color4f_t* color4f);
SK_C_API void sk_color4f_from_color(sk_color_t color, sk_color4f_t* color4f);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_data.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_data_DEFINED
#define sk_data_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API sk_data_t* sk_data_new_empty(void);
SK_C_API sk_data_t* sk_data_new_with_copy(const void* src, size_t length);
SK_C_API sk_data_t* sk_data_new_subset(const sk_data_t* src, size_t offset, size_t length);
SK_C_API void sk_data_ref(const sk_data_t*);
SK_C_API void sk_data_unref(const sk_data_t*);
SK_C_API size_t sk_data_get_size(const sk_data_t*);
SK_C_API const void* sk_data_get_data(const sk_data_t*);
SK_C_API sk_data_t* sk_data_new_from_file(const char* path);
SK_C_API sk_data_t* sk_data_new_from_stream(sk_stream_t* stream, size_t length);
SK_C_API const uint8_t* sk_data_get_bytes(const sk_data_t*);
SK_C_API sk_data_t* sk_data_new_with_proc(const void* ptr, size_t length, sk_data_release_proc proc, void* ctx);
SK_C_API sk_data_t* sk_data_new_uninitialized(size_t size);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_document.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_document_DEFINED
#define sk_document_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void sk_document_unref(sk_document_t* document);

SK_C_API sk_document_t* sk_document_create_pdf_from_stream(sk_wstream_t* stream);
SK_C_API sk_document_t* sk_document_create_pdf_from_stream_with_metadata(sk_wstream_t* stream, const sk_document_pdf_metadata_t* metadata);

SK_C_API sk_document_t* sk_document_create_xps_from_stream(sk_wstream_t* stream, float dpi);

SK_C_API sk_canvas_t* sk_document_begin_page(sk_document_t* document, float width, float height, const sk_rect_t* content);
SK_C_API void sk_document_end_page(sk_document_t* document);
SK_C_API void sk_document_close(sk_document_t* document);
SK_C_API void sk_document_abort(sk_document_t* document);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_drawable.h
/*
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_drawable_DEFINED
#define sk_drawable_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

typedef struct sk_drawable_t sk_drawable_t;

SK_C_API void sk_drawable_unref (sk_drawable_t*);
SK_C_API uint32_t sk_drawable_get_generation_id (sk_drawable_t*);
SK_C_API void sk_drawable_get_bounds (sk_drawable_t*, sk_rect_t*);
SK_C_API void sk_drawable_draw (sk_drawable_t*, sk_canvas_t*, const sk_matrix_t*);
SK_C_API sk_picture_t* sk_drawable_new_picture_snapshot(sk_drawable_t*);
SK_C_API void sk_drawable_notify_drawing_changed (sk_drawable_t*);
SK_C_API size_t sk_drawable_approximate_bytes_used(sk_drawable_t*);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_font.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_font_DEFINED
#define sk_font_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

// sk_font_t

SK_C_API sk_font_t* sk_font_new(void);
SK_C_API sk_font_t* sk_font_new_with_values(sk_typeface_t* typeface, float size, float scaleX, float skewX);
SK_C_API void sk_font_delete(sk_font_t* font);
SK_C_API bool sk_font_is_force_auto_hinting(const sk_font_t* font);
SK_C_API void sk_font_set_force_auto_hinting(sk_font_t* font, bool value);
SK_C_API bool sk_font_is_embedded_bitmaps(const sk_font_t* font);
SK_C_API void sk_font_set_embedded_bitmaps(sk_font_t* font, bool value);
SK_C_API bool sk_font_is_subpixel(const sk_font_t* font);
SK_C_API void sk_font_set_subpixel(sk_font_t* font, bool value);
SK_C_API bool sk_font_is_linear_metrics(const sk_font_t* font);
SK_C_API void sk_font_set_linear_metrics(sk_font_t* font, bool value);
SK_C_API bool sk_font_is_embolden(const sk_font_t* font);
SK_C_API void sk_font_set_embolden(sk_font_t* font, bool value);
SK_C_API bool sk_font_is_baseline_snap(const sk_font_t* font);
SK_C_API void sk_font_set_baseline_snap(sk_font_t* font, bool value);
SK_C_API sk_font_edging_t sk_font_get_edging(const sk_font_t* font);
SK_C_API void sk_font_set_edging(sk_font_t* font, sk_font_edging_t value);
SK_C_API sk_font_hinting_t sk_font_get_hinting(const sk_font_t* font);
SK_C_API void sk_font_set_hinting(sk_font_t* font, sk_font_hinting_t value);
SK_C_API sk_typeface_t* sk_font_get_typeface(const sk_font_t* font);
SK_C_API void sk_font_set_typeface(sk_font_t* font, sk_typeface_t* value);
SK_C_API float sk_font_get_size(const sk_font_t* font);
SK_C_API void sk_font_set_size(sk_font_t* font, float value);
SK_C_API float sk_font_get_scale_x(const sk_font_t* font);
SK_C_API void sk_font_set_scale_x(sk_font_t* font, float value);
SK_C_API float sk_font_get_skew_x(const sk_font_t* font);
SK_C_API void sk_font_set_skew_x(sk_font_t* font, float value);
SK_C_API int sk_font_text_to_glyphs(const sk_font_t* font, const void* text, size_t byteLength, sk_text_encoding_t encoding, uint16_t glyphs[], int maxGlyphCount);
SK_C_API uint16_t sk_font_unichar_to_glyph(const sk_font_t* font, int32_t uni);
SK_C_API void sk_font_unichars_to_glyphs(const sk_font_t* font, const int32_t uni[], int count, uint16_t glyphs[]);
SK_C_API float sk_font_measure_text(const sk_font_t* font, const void* text, size_t byteLength, sk_text_encoding_t encoding, sk_rect_t* bounds, const sk_paint_t* paint);
// NOTE: it appears that .NET Framework 4.7 has an issue with returning float?
//       https://github.com/mono/SkiaSharp/issues/1409
SK_C_API void sk_font_measure_text_no_return(const sk_font_t* font, const void* text, size_t byteLength, sk_text_encoding_t encoding, sk_rect_t* bounds, const sk_paint_t* paint, float* measuredWidth);
SK_C_API size_t sk_font_break_text(const sk_font_t* font, const void* text, size_t byteLength, sk_text_encoding_t encoding, float maxWidth, float* measuredWidth, const sk_paint_t* paint);
SK_C_API void sk_font_get_widths_bounds(const sk_font_t* font, const uint16_t glyphs[], int count, float widths[], sk_rect_t bounds[], const sk_paint_t* paint);
SK_C_API void sk_font_get_pos(const sk_font_t* font, const uint16_t glyphs[], int count, sk_point_t pos[], sk_point_t* origin);
SK_C_API void sk_font_get_xpos(const sk_font_t* font, const uint16_t glyphs[], int count, float xpos[], float origin);
SK_C_API bool sk_font_get_path(const sk_font_t* font, uint16_t glyph, sk_path_t* path);
SK_C_API void sk_font_get_paths(const sk_font_t* font, uint16_t glyphs[], int count, const sk_glyph_path_proc glyphPathProc, void* context);
SK_C_API float sk_font_get_metrics(const sk_font_t* font, sk_fontmetrics_t* metrics);

// sk_text_utils

SK_C_API void sk_text_utils_get_path(const void* text, size_t length, sk_text_encoding_t encoding, float x, float y, const sk_font_t* font, sk_path_t* path);
SK_C_API void sk_text_utils_get_pos_path(const void* text, size_t length, sk_text_encoding_t encoding, const sk_point_t pos[], const sk_font_t* font, sk_path_t* path);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_general.h
/*
 * Copyright 2019 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_general_DEFINED
#define sk_general_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

// ref counting

SK_C_API bool sk_refcnt_unique(const sk_refcnt_t* refcnt);
SK_C_API int sk_refcnt_get_ref_count(const sk_refcnt_t* refcnt);
SK_C_API void sk_refcnt_safe_ref(sk_refcnt_t* refcnt);
SK_C_API void sk_refcnt_safe_unref(sk_refcnt_t* refcnt);

SK_C_API bool sk_nvrefcnt_unique(const sk_nvrefcnt_t* refcnt);
SK_C_API int sk_nvrefcnt_get_ref_count(const sk_nvrefcnt_t* refcnt);
SK_C_API void sk_nvrefcnt_safe_ref(sk_nvrefcnt_t* refcnt);
SK_C_API void sk_nvrefcnt_safe_unref(sk_nvrefcnt_t* refcnt);

// color type

SK_C_API sk_colortype_t sk_colortype_get_default_8888(void);

// library information

SK_C_API int sk_version_get_milestone(void);
SK_C_API int sk_version_get_increment(void);
SK_C_API const char* sk_version_get_string(void);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_graphics.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_graphics_DEFINED
#define sk_graphics_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD


SK_C_API void sk_graphics_init(void);

// purge
SK_C_API void sk_graphics_purge_font_cache(void);
SK_C_API void sk_graphics_purge_resource_cache(void);
SK_C_API void sk_graphics_purge_all_caches(void);

// font cache
SK_C_API size_t sk_graphics_get_font_cache_used(void);
SK_C_API size_t sk_graphics_get_font_cache_limit(void);
SK_C_API size_t sk_graphics_set_font_cache_limit(size_t bytes);
SK_C_API int sk_graphics_get_font_cache_count_used(void);
SK_C_API int sk_graphics_get_font_cache_count_limit(void);
SK_C_API int sk_graphics_set_font_cache_count_limit(int count);

// resource cache
SK_C_API size_t sk_graphics_get_resource_cache_total_bytes_used(void);
SK_C_API size_t sk_graphics_get_resource_cache_total_byte_limit(void);
SK_C_API size_t sk_graphics_set_resource_cache_total_byte_limit(size_t newLimit);
SK_C_API size_t sk_graphics_get_resource_cache_single_allocation_byte_limit(void);
SK_C_API size_t sk_graphics_set_resource_cache_single_allocation_byte_limit(size_t newLimit);

// dump
SK_C_API void sk_graphics_dump_memory_statistics(sk_tracememorydump_t* dump);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_image.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_image_DEFINED
#define sk_image_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void sk_image_ref(const sk_image_t* cimage);
SK_C_API void sk_image_unref(const sk_image_t* cimage);
SK_C_API sk_image_t* sk_image_new_raster_copy(const sk_imageinfo_t* cinfo, const void* pixels, size_t rowBytes);
SK_C_API sk_image_t* sk_image_new_raster_copy_with_pixmap(const sk_pixmap_t* pixmap);
SK_C_API sk_image_t* sk_image_new_raster_data(const sk_imageinfo_t* cinfo, sk_data_t* pixels, size_t rowBytes);
SK_C_API sk_image_t* sk_image_new_raster(const sk_pixmap_t* pixmap, sk_image_raster_release_proc releaseProc, void* context);
SK_C_API sk_image_t* sk_image_new_from_bitmap(const sk_bitmap_t* cbitmap);
SK_C_API sk_image_t* sk_image_new_from_encoded(const sk_data_t* cdata);
SK_C_API sk_image_t* sk_image_new_from_texture(gr_recording_context_t* context, const gr_backendtexture_t* texture, gr_surfaceorigin_t origin, sk_colortype_t colorType, sk_alphatype_t alpha, const sk_colorspace_t* colorSpace, const sk_image_texture_release_proc releaseProc, void* releaseContext);
SK_C_API sk_image_t* sk_image_new_from_adopted_texture(gr_recording_context_t* context, const gr_backendtexture_t* texture, gr_surfaceorigin_t origin, sk_colortype_t colorType, sk_alphatype_t alpha, const sk_colorspace_t* colorSpace);
SK_C_API sk_image_t* sk_image_new_from_picture(sk_picture_t* picture, const sk_isize_t* dimensions, const sk_matrix_t* cmatrix, const sk_paint_t* paint, bool useFloatingPointBitDepth, const sk_colorspace_t* colorSpace, const sk_surfaceprops_t* props);
SK_C_API int sk_image_get_width(const sk_image_t* cimage);
SK_C_API int sk_image_get_height(const sk_image_t* cimage);
SK_C_API uint32_t sk_image_get_unique_id(const sk_image_t* cimage);
SK_C_API sk_alphatype_t sk_image_get_alpha_type(const sk_image_t* image);
SK_C_API sk_colortype_t sk_image_get_color_type(const sk_image_t* image);
SK_C_API sk_colorspace_t* sk_image_get_colorspace(const sk_image_t* image);
SK_C_API bool sk_image_is_alpha_only(const sk_image_t* image);
SK_C_API sk_shader_t* sk_image_make_shader(const sk_image_t* image, sk_shader_tilemode_t tileX, sk_shader_tilemode_t tileY, const sk_sampling_options_t* sampling, const sk_matrix_t* cmatrix);
SK_C_API sk_shader_t* sk_image_make_raw_shader(const sk_image_t* image, sk_shader_tilemode_t tileX, sk_shader_tilemode_t tileY, const sk_sampling_options_t* sampling, const sk_matrix_t* cmatrix);
SK_C_API bool sk_image_peek_pixels(const sk_image_t* image, sk_pixmap_t* pixmap);
SK_C_API bool sk_image_is_texture_backed(const sk_image_t* image);
SK_C_API bool sk_image_is_lazy_generated(const sk_image_t* image);
SK_C_API bool sk_image_is_valid(const sk_image_t* image, gr_recording_context_t* context);
SK_C_API bool sk_image_read_pixels(const sk_image_t* image, const sk_imageinfo_t* dstInfo, void* dstPixels, size_t dstRowBytes, int srcX, int srcY, sk_image_caching_hint_t cachingHint);
SK_C_API bool sk_image_read_pixels_into_pixmap(const sk_image_t* image, const sk_pixmap_t* dst, int srcX, int srcY, sk_image_caching_hint_t cachingHint);
SK_C_API bool sk_image_scale_pixels(const sk_image_t* image, const sk_pixmap_t* dst, const sk_sampling_options_t* sampling, sk_image_caching_hint_t cachingHint);
SK_C_API sk_data_t* sk_image_ref_encoded(const sk_image_t* cimage);
SK_C_API sk_image_t* sk_image_make_subset_raster(const sk_image_t* cimage, const sk_irect_t* subset);
SK_C_API sk_image_t* sk_image_make_subset(const sk_image_t* cimage, gr_direct_context_t* context, const sk_irect_t* subset);
SK_C_API sk_image_t* sk_image_make_texture_image(const sk_image_t* cimage, gr_direct_context_t* context, bool mipmapped, bool budgeted);
SK_C_API sk_image_t* sk_image_make_non_texture_image(const sk_image_t* cimage);
SK_C_API sk_image_t* sk_image_make_raster_image(const sk_image_t* cimage);
SK_C_API sk_image_t* sk_image_make_with_filter_raster(const sk_image_t* cimage, const sk_imagefilter_t* filter, const sk_irect_t* subset, const sk_irect_t* clipBounds, sk_irect_t* outSubset, sk_ipoint_t* outOffset);
SK_C_API sk_image_t* sk_image_make_with_filter(const sk_image_t* cimage, gr_recording_context_t* context, const sk_imagefilter_t* filter, const sk_irect_t* subset, const sk_irect_t* clipBounds, sk_irect_t* outSubset, sk_ipoint_t* outOffset);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_imagefilter.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_imagefilter_DEFINED
#define sk_imagefilter_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD


// sk_imagefilter_t

SK_C_API void sk_imagefilter_unref(sk_imagefilter_t* cfilter);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_arithmetic(float k1, float k2, float k3, float k4, bool enforcePMColor, const sk_imagefilter_t* background, const sk_imagefilter_t* foreground, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_blend(sk_blendmode_t mode, const sk_imagefilter_t* background, const sk_imagefilter_t* foreground, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_blender(sk_blender_t* blender, const sk_imagefilter_t* background, const sk_imagefilter_t* foreground, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_blur(float sigmaX, float sigmaY, sk_shader_tilemode_t tileMode, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_color_filter(sk_colorfilter_t* cf, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_compose(const sk_imagefilter_t* outer, const sk_imagefilter_t* inner);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_displacement_map_effect(sk_color_channel_t xChannelSelector, sk_color_channel_t yChannelSelector, float scale, const sk_imagefilter_t* displacement, const sk_imagefilter_t* color, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_drop_shadow(float dx, float dy, float sigmaX, float sigmaY, sk_color_t color, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_drop_shadow_only(float dx, float dy, float sigmaX, float sigmaY, sk_color_t color, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_image(sk_image_t* image, const sk_rect_t* srcRect, const sk_rect_t* dstRect, const sk_sampling_options_t* sampling);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_image_simple(sk_image_t* image, const sk_sampling_options_t* sampling);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_magnifier(const sk_rect_t* lensBounds, float zoomAmount, float inset, const sk_sampling_options_t* sampling, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_matrix_convolution(const sk_isize_t* kernelSize, const float kernel[], float gain, float bias, const sk_ipoint_t* kernelOffset, sk_shader_tilemode_t ctileMode, bool convolveAlpha, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_matrix_transform(const sk_matrix_t* cmatrix, const sk_sampling_options_t* sampling, const sk_imagefilter_t* input);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_merge(const sk_imagefilter_t* cfilters[], int count, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_merge_simple(const sk_imagefilter_t* first, const sk_imagefilter_t* second, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_offset(float dx, float dy, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_picture(const sk_picture_t* picture);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_picture_with_rect(const sk_picture_t* picture, const sk_rect_t* targetRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_shader(const sk_shader_t* shader, bool dither, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_tile(const sk_rect_t* src, const sk_rect_t* dst, const sk_imagefilter_t* input);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_dilate(float radiusX, float radiusY, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_erode(float radiusX, float radiusY, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_distant_lit_diffuse(const sk_point3_t* direction, sk_color_t lightColor, float surfaceScale, float kd, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_point_lit_diffuse(const sk_point3_t* location, sk_color_t lightColor, float surfaceScale, float kd, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_spot_lit_diffuse(const sk_point3_t* location, const sk_point3_t* target, float specularExponent, float cutoffAngle, sk_color_t lightColor, float surfaceScale, float kd, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_distant_lit_specular(const sk_point3_t* direction, sk_color_t lightColor, float surfaceScale, float ks, float shininess, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_point_lit_specular(const sk_point3_t* location, sk_color_t lightColor, float surfaceScale, float ks, float shininess, const sk_imagefilter_t* input, const sk_rect_t* cropRect);
SK_C_API sk_imagefilter_t* sk_imagefilter_new_spot_lit_specular(const sk_point3_t* location, const sk_point3_t* target, float specularExponent, float cutoffAngle, sk_color_t lightColor, float surfaceScale, float ks, float shininess, const sk_imagefilter_t* input, const sk_rect_t* cropRect);


SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_linker.h
/*
 * Copyright 2024 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_linker_DEFINED
#define sk_linker_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void sk_linker_keep_alive(void);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_maskfilter.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_maskfilter_DEFINED
#define sk_maskfilter_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void sk_maskfilter_ref(sk_maskfilter_t*);
SK_C_API void sk_maskfilter_unref(sk_maskfilter_t*);
SK_C_API sk_maskfilter_t* sk_maskfilter_new_blur(sk_blurstyle_t, float sigma);
SK_C_API sk_maskfilter_t* sk_maskfilter_new_blur_with_flags(sk_blurstyle_t, float sigma, bool respectCTM);
SK_C_API sk_maskfilter_t* sk_maskfilter_new_table(const uint8_t table[256]);
SK_C_API sk_maskfilter_t* sk_maskfilter_new_gamma(float gamma);
SK_C_API sk_maskfilter_t* sk_maskfilter_new_clip(uint8_t min, uint8_t max);
SK_C_API sk_maskfilter_t* sk_maskfilter_new_shader(sk_shader_t* cshader);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_matrix.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_matrix_DEFINED
#define sk_matrix_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API bool sk_matrix_try_invert (sk_matrix_t *matrix, sk_matrix_t *result);
SK_C_API void sk_matrix_concat (sk_matrix_t *result, sk_matrix_t *first, sk_matrix_t *second);
SK_C_API void sk_matrix_pre_concat (sk_matrix_t *result, sk_matrix_t *matrix);
SK_C_API void sk_matrix_post_concat (sk_matrix_t *result, sk_matrix_t *matrix);
SK_C_API void sk_matrix_map_rect (sk_matrix_t *matrix, sk_rect_t *dest, sk_rect_t *source);
SK_C_API void sk_matrix_map_points (sk_matrix_t *matrix, sk_point_t *dst, sk_point_t *src, int count);
SK_C_API void sk_matrix_map_vectors (sk_matrix_t *matrix, sk_point_t *dst, sk_point_t *src, int count);
SK_C_API void sk_matrix_map_xy (sk_matrix_t *matrix, float x, float y, sk_point_t* result);
SK_C_API void sk_matrix_map_vector (sk_matrix_t *matrix, float x, float y, sk_point_t* result);
SK_C_API float sk_matrix_map_radius (sk_matrix_t *matrix, float radius);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_paint.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_paint_DEFINED
#define sk_paint_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API sk_paint_t* sk_paint_new(void);
SK_C_API sk_paint_t* sk_paint_clone(sk_paint_t*);
SK_C_API void sk_paint_delete(sk_paint_t*);
SK_C_API void sk_paint_reset(sk_paint_t*);
SK_C_API bool sk_paint_is_antialias(const sk_paint_t*);
SK_C_API void sk_paint_set_antialias(sk_paint_t*, bool);
SK_C_API sk_color_t sk_paint_get_color(const sk_paint_t*);
SK_C_API void sk_paint_get_color4f(const sk_paint_t* paint, sk_color4f_t* color);
SK_C_API void sk_paint_set_color(sk_paint_t*, sk_color_t);
SK_C_API void sk_paint_set_color4f(sk_paint_t* paint, sk_color4f_t* color, sk_colorspace_t* colorspace);
SK_C_API sk_paint_style_t sk_paint_get_style(const sk_paint_t*);
SK_C_API void sk_paint_set_style(sk_paint_t*, sk_paint_style_t);
SK_C_API float sk_paint_get_stroke_width(const sk_paint_t*);
SK_C_API void sk_paint_set_stroke_width(sk_paint_t*, float width);
SK_C_API float sk_paint_get_stroke_miter(const sk_paint_t*);
SK_C_API void sk_paint_set_stroke_miter(sk_paint_t*, float miter);
SK_C_API sk_stroke_cap_t sk_paint_get_stroke_cap(const sk_paint_t*);
SK_C_API void sk_paint_set_stroke_cap(sk_paint_t*, sk_stroke_cap_t);
SK_C_API sk_stroke_join_t sk_paint_get_stroke_join(const sk_paint_t*);
SK_C_API void sk_paint_set_stroke_join(sk_paint_t*, sk_stroke_join_t);
SK_C_API void sk_paint_set_shader(sk_paint_t*, sk_shader_t*);
SK_C_API void sk_paint_set_maskfilter(sk_paint_t*, sk_maskfilter_t*);
SK_C_API void sk_paint_set_blendmode(sk_paint_t*, sk_blendmode_t);
SK_C_API void sk_paint_set_blender(sk_paint_t* paint, sk_blender_t* blender);
SK_C_API bool sk_paint_is_dither(const sk_paint_t*);
SK_C_API void sk_paint_set_dither(sk_paint_t*, bool);
SK_C_API sk_shader_t* sk_paint_get_shader(sk_paint_t*);
SK_C_API sk_maskfilter_t* sk_paint_get_maskfilter(sk_paint_t*);
SK_C_API void sk_paint_set_colorfilter(sk_paint_t*, sk_colorfilter_t*);
SK_C_API sk_colorfilter_t* sk_paint_get_colorfilter(sk_paint_t*);
SK_C_API void sk_paint_set_imagefilter(sk_paint_t*, sk_imagefilter_t*);
SK_C_API sk_imagefilter_t* sk_paint_get_imagefilter(sk_paint_t*);
SK_C_API sk_blendmode_t sk_paint_get_blendmode(sk_paint_t*);
SK_C_API sk_blender_t* sk_paint_get_blender(sk_paint_t* cpaint);
SK_C_API sk_path_effect_t* sk_paint_get_path_effect(sk_paint_t* cpaint);
SK_C_API void sk_paint_set_path_effect(sk_paint_t* cpaint, sk_path_effect_t* effect);  
SK_C_API bool sk_paint_get_fill_path(const sk_paint_t* cpaint, const sk_path_t* src, sk_path_t* dst, const sk_rect_t* cullRect, const sk_matrix_t* cmatrix);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_path.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_path_DEFINED
#define sk_path_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

/* Path */
SK_C_API sk_path_t* sk_path_new(void);
SK_C_API void sk_path_delete(sk_path_t*);
SK_C_API void sk_path_move_to(sk_path_t*, float x, float y);
SK_C_API void sk_path_line_to(sk_path_t*, float x, float y);
SK_C_API void sk_path_quad_to(sk_path_t*, float x0, float y0, float x1, float y1);
SK_C_API void sk_path_conic_to(sk_path_t*, float x0, float y0, float x1, float y1, float w);
SK_C_API void sk_path_cubic_to(sk_path_t*, float x0, float y0, float x1, float y1, float x2, float y2);
SK_C_API void sk_path_arc_to(sk_path_t*, float rx, float ry, float xAxisRotate, sk_path_arc_size_t largeArc, sk_path_direction_t sweep, float x, float y);
SK_C_API void sk_path_rarc_to(sk_path_t*, float rx, float ry, float xAxisRotate, sk_path_arc_size_t largeArc, sk_path_direction_t sweep, float x, float y);
SK_C_API void sk_path_arc_to_with_oval(sk_path_t*, const sk_rect_t* oval, float startAngle, float sweepAngle, bool forceMoveTo);
SK_C_API void sk_path_arc_to_with_points(sk_path_t*, float x1, float y1, float x2, float y2, float radius);
SK_C_API void sk_path_close(sk_path_t*);
SK_C_API void sk_path_add_rect(sk_path_t*, const sk_rect_t*, sk_path_direction_t);
SK_C_API void sk_path_add_rrect(sk_path_t*, const sk_rrect_t*, sk_path_direction_t);
SK_C_API void sk_path_add_rrect_start(sk_path_t*, const sk_rrect_t*, sk_path_direction_t, uint32_t);
SK_C_API void sk_path_add_rounded_rect(sk_path_t*, const sk_rect_t*, float, float, sk_path_direction_t);
SK_C_API void sk_path_add_oval(sk_path_t*, const sk_rect_t*, sk_path_direction_t);
SK_C_API void sk_path_add_circle(sk_path_t*, float x, float y, float radius, sk_path_direction_t dir);
SK_C_API void sk_path_get_bounds(const sk_path_t*, sk_rect_t*);
SK_C_API void sk_path_compute_tight_bounds(const sk_path_t*, sk_rect_t*);
SK_C_API void sk_path_rmove_to(sk_path_t*, float dx, float dy);
SK_C_API void sk_path_rline_to(sk_path_t*, float dx, float yd);
SK_C_API void sk_path_rquad_to(sk_path_t*, float dx0, float dy0, float dx1, float dy1);
SK_C_API void sk_path_rconic_to(sk_path_t*, float dx0, float dy0, float dx1, float dy1, float w);
SK_C_API void sk_path_rcubic_to(sk_path_t*, float dx0, float dy0, float dx1, float dy1, float dx2, float dy2);
SK_C_API void sk_path_add_rect_start(sk_path_t* cpath, const sk_rect_t* crect, sk_path_direction_t cdir, uint32_t startIndex);
SK_C_API void sk_path_add_arc(sk_path_t* cpath, const sk_rect_t* crect, float startAngle, float sweepAngle);
SK_C_API sk_path_filltype_t sk_path_get_filltype(sk_path_t*);
SK_C_API void sk_path_set_filltype(sk_path_t*, sk_path_filltype_t);
SK_C_API void sk_path_transform(sk_path_t* cpath, const sk_matrix_t* cmatrix);
SK_C_API void sk_path_transform_to_dest(const sk_path_t* cpath, const sk_matrix_t* cmatrix, sk_path_t* destination);
SK_C_API sk_path_t* sk_path_clone(const sk_path_t* cpath);
SK_C_API void sk_path_add_path_offset  (sk_path_t* cpath, sk_path_t* other, float dx, float dy, sk_path_add_mode_t add_mode);
SK_C_API void sk_path_add_path_matrix  (sk_path_t* cpath, sk_path_t* other, sk_matrix_t *matrix, sk_path_add_mode_t add_mode);
SK_C_API void sk_path_add_path         (sk_path_t* cpath, sk_path_t* other, sk_path_add_mode_t add_mode);
SK_C_API void sk_path_add_path_reverse (sk_path_t* cpath, sk_path_t* other);
SK_C_API void sk_path_reset (sk_path_t* cpath);
SK_C_API void sk_path_rewind (sk_path_t* cpath);
SK_C_API int sk_path_count_points (const sk_path_t* cpath);
SK_C_API int sk_path_count_verbs (const sk_path_t* cpath);
SK_C_API void sk_path_get_point (const sk_path_t* cpath, int index, sk_point_t* point);
SK_C_API int sk_path_get_points (const sk_path_t* cpath, sk_point_t* points, int max);
SK_C_API bool sk_path_contains (const sk_path_t* cpath, float x, float y);
SK_C_API bool sk_path_parse_svg_string (sk_path_t* cpath, const char* str);
SK_C_API void sk_path_to_svg_string (const sk_path_t* cpath, sk_string_t* str);
SK_C_API bool sk_path_get_last_point (const sk_path_t* cpath, sk_point_t* point);
SK_C_API int sk_path_convert_conic_to_quads(const sk_point_t* p0, const sk_point_t* p1, const sk_point_t* p2, float w, sk_point_t* pts, int pow2);
SK_C_API void sk_path_add_poly(sk_path_t* cpath, const sk_point_t* points, int count, bool close);
SK_C_API uint32_t sk_path_get_segment_masks(sk_path_t* cpath);
SK_C_API bool sk_path_is_oval(sk_path_t* cpath, sk_rect_t* bounds);
SK_C_API bool sk_path_is_rrect(sk_path_t* cpath, sk_rrect_t* bounds);
SK_C_API bool sk_path_is_line(sk_path_t* cpath, sk_point_t line [2]);
SK_C_API bool sk_path_is_rect(sk_path_t* cpath, sk_rect_t* rect, bool* isClosed, sk_path_direction_t* direction);
SK_C_API bool sk_path_is_convex(const sk_path_t* cpath);

/* Iterators */
SK_C_API sk_path_iterator_t* sk_path_create_iter (sk_path_t *cpath, int forceClose);
SK_C_API sk_path_verb_t sk_path_iter_next (sk_path_iterator_t *iterator, sk_point_t points [4]);
SK_C_API float sk_path_iter_conic_weight (sk_path_iterator_t *iterator);
SK_C_API int sk_path_iter_is_close_line (sk_path_iterator_t *iterator);
SK_C_API int sk_path_iter_is_closed_contour (sk_path_iterator_t *iterator);
SK_C_API void sk_path_iter_destroy (sk_path_iterator_t *iterator);

/* Raw iterators */
SK_C_API sk_path_rawiterator_t* sk_path_create_rawiter (sk_path_t *cpath);
SK_C_API sk_path_verb_t sk_path_rawiter_peek (sk_path_rawiterator_t *iterator);
SK_C_API sk_path_verb_t sk_path_rawiter_next (sk_path_rawiterator_t *iterator, sk_point_t points [4]);
SK_C_API float sk_path_rawiter_conic_weight (sk_path_rawiterator_t *iterator);
SK_C_API void sk_path_rawiter_destroy (sk_path_rawiterator_t *iterator);

/* Path Ops */
SK_C_API bool sk_pathop_op(const sk_path_t* one, const sk_path_t* two, sk_pathop_t op, sk_path_t* result);
SK_C_API bool sk_pathop_simplify(const sk_path_t* path, sk_path_t* result);
SK_C_API bool sk_pathop_tight_bounds(const sk_path_t* path, sk_rect_t* result);
SK_C_API bool sk_pathop_as_winding(const sk_path_t* path, sk_path_t* result);

/* Path Op Builder */
SK_C_API sk_opbuilder_t* sk_opbuilder_new(void);
SK_C_API void sk_opbuilder_destroy(sk_opbuilder_t* builder);
SK_C_API void sk_opbuilder_add(sk_opbuilder_t* builder, const sk_path_t* path, sk_pathop_t op);
SK_C_API bool sk_opbuilder_resolve(sk_opbuilder_t* builder, sk_path_t* result);

/* Path Measure */
SK_C_API sk_pathmeasure_t* sk_pathmeasure_new(void);
SK_C_API sk_pathmeasure_t* sk_pathmeasure_new_with_path(const sk_path_t* path, bool forceClosed, float resScale);
SK_C_API void sk_pathmeasure_destroy(sk_pathmeasure_t* pathMeasure);
SK_C_API void sk_pathmeasure_set_path(sk_pathmeasure_t* pathMeasure, const sk_path_t* path, bool forceClosed);
SK_C_API float sk_pathmeasure_get_length(sk_pathmeasure_t* pathMeasure);
SK_C_API bool sk_pathmeasure_get_pos_tan(sk_pathmeasure_t* pathMeasure, float distance, sk_point_t* position, sk_vector_t* tangent);
SK_C_API bool sk_pathmeasure_get_matrix(sk_pathmeasure_t* pathMeasure, float distance, sk_matrix_t* matrix, sk_pathmeasure_matrixflags_t flags);
SK_C_API bool sk_pathmeasure_get_segment(sk_pathmeasure_t* pathMeasure, float start, float stop, sk_path_t* dst, bool startWithMoveTo);
SK_C_API bool sk_pathmeasure_is_closed(sk_pathmeasure_t* pathMeasure);
SK_C_API bool sk_pathmeasure_next_contour(sk_pathmeasure_t* pathMeasure);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_patheffect.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_patheffect_DEFINED
#define sk_patheffect_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void sk_path_effect_unref(sk_path_effect_t* t); 
SK_C_API sk_path_effect_t* sk_path_effect_create_compose(sk_path_effect_t* outer, sk_path_effect_t* inner);
SK_C_API sk_path_effect_t* sk_path_effect_create_sum(sk_path_effect_t* first, sk_path_effect_t* second);
SK_C_API sk_path_effect_t* sk_path_effect_create_discrete(float segLength, float deviation, uint32_t seedAssist /*0*/);
SK_C_API sk_path_effect_t* sk_path_effect_create_corner(float radius);
SK_C_API sk_path_effect_t* sk_path_effect_create_1d_path(const sk_path_t* path, float advance, float phase, sk_path_effect_1d_style_t style);
SK_C_API sk_path_effect_t* sk_path_effect_create_2d_line(float width, const sk_matrix_t* matrix);
SK_C_API sk_path_effect_t* sk_path_effect_create_2d_path(const sk_matrix_t* matrix, const sk_path_t* path);
SK_C_API sk_path_effect_t* sk_path_effect_create_dash(const float intervals[], int count, float phase);
SK_C_API sk_path_effect_t* sk_path_effect_create_trim(float start, float stop, sk_path_effect_trim_mode_t mode);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_picture.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_picture_DEFINED
#define sk_picture_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

// SkPictureRecorder

SK_C_API sk_picture_recorder_t* sk_picture_recorder_new(void);
SK_C_API void sk_picture_recorder_delete(sk_picture_recorder_t*);
SK_C_API sk_canvas_t* sk_picture_recorder_begin_recording(sk_picture_recorder_t*, const sk_rect_t*);
SK_C_API sk_canvas_t* sk_picture_recorder_begin_recording_with_bbh_factory(sk_picture_recorder_t*, const sk_rect_t*, sk_bbh_factory_t*);
SK_C_API sk_picture_t* sk_picture_recorder_end_recording(sk_picture_recorder_t*);
SK_C_API sk_drawable_t* sk_picture_recorder_end_recording_as_drawable(sk_picture_recorder_t*);
SK_C_API sk_canvas_t* sk_picture_get_recording_canvas(sk_picture_recorder_t* crec);

// SkPicture

SK_C_API void sk_picture_ref(sk_picture_t*);
SK_C_API void sk_picture_unref(sk_picture_t*);
SK_C_API uint32_t sk_picture_get_unique_id(sk_picture_t*);
SK_C_API void sk_picture_get_cull_rect(sk_picture_t*, sk_rect_t*);
SK_C_API sk_shader_t* sk_picture_make_shader(sk_picture_t* src, sk_shader_tilemode_t tmx, sk_shader_tilemode_t tmy, sk_filter_mode_t mode, const sk_matrix_t* localMatrix, const sk_rect_t* tile);
SK_C_API sk_data_t* sk_picture_serialize_to_data(const sk_picture_t* picture);
SK_C_API void sk_picture_serialize_to_stream(const sk_picture_t* picture, sk_wstream_t* stream);
SK_C_API sk_picture_t* sk_picture_deserialize_from_stream(sk_stream_t* stream);
SK_C_API sk_picture_t* sk_picture_deserialize_from_data(sk_data_t* data);
SK_C_API sk_picture_t* sk_picture_deserialize_from_memory(void* buffer, size_t length);
SK_C_API void sk_picture_playback(const sk_picture_t* picture, sk_canvas_t* canvas);
SK_C_API int sk_picture_approximate_op_count(const sk_picture_t* picture, bool nested);
SK_C_API size_t sk_picture_approximate_bytes_used(const sk_picture_t* picture);

// SkRTreeFactory

SK_C_API sk_rtree_factory_t* sk_rtree_factory_new(void);
SK_C_API void sk_rtree_factory_delete(sk_rtree_factory_t*);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_pixmap.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_pixmap_DEFINED
#define sk_pixmap_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

// SkPixmap

SK_C_API void sk_pixmap_destructor(sk_pixmap_t* cpixmap);
SK_C_API sk_pixmap_t* sk_pixmap_new(void);
SK_C_API sk_pixmap_t* sk_pixmap_new_with_params(const sk_imageinfo_t* cinfo, const void* addr, size_t rowBytes);
SK_C_API void sk_pixmap_reset(sk_pixmap_t* cpixmap);
SK_C_API void sk_pixmap_reset_with_params(sk_pixmap_t* cpixmap, const sk_imageinfo_t* cinfo, const void* addr, size_t rowBytes);
SK_C_API void sk_pixmap_set_colorspace(sk_pixmap_t* cpixmap, sk_colorspace_t* colorspace);
SK_C_API bool sk_pixmap_extract_subset(const sk_pixmap_t* cpixmap, sk_pixmap_t* result, const sk_irect_t* subset);
SK_C_API void sk_pixmap_get_info(const sk_pixmap_t* cpixmap, sk_imageinfo_t* cinfo);
SK_C_API size_t sk_pixmap_get_row_bytes(const sk_pixmap_t* cpixmap);
SK_C_API sk_colorspace_t* sk_pixmap_get_colorspace(const sk_pixmap_t* cpixmap);
SK_C_API bool sk_pixmap_compute_is_opaque(const sk_pixmap_t* cpixmap);
SK_C_API sk_color_t sk_pixmap_get_pixel_color(const sk_pixmap_t* cpixmap, int x, int y);
SK_C_API void sk_pixmap_get_pixel_color4f(const sk_pixmap_t* cpixmap, int x, int y, sk_color4f_t* color);
SK_C_API float sk_pixmap_get_pixel_alphaf(const sk_pixmap_t* cpixmap, int x, int y);
SK_C_API void* sk_pixmap_get_writable_addr(const sk_pixmap_t* cpixmap);
SK_C_API void* sk_pixmap_get_writeable_addr_with_xy(const sk_pixmap_t* cpixmap, int x, int y);
SK_C_API bool sk_pixmap_read_pixels(const sk_pixmap_t* cpixmap, const sk_imageinfo_t* dstInfo, void* dstPixels, size_t dstRowBytes, int srcX, int srcY);
SK_C_API bool sk_pixmap_scale_pixels(const sk_pixmap_t* cpixmap, const sk_pixmap_t* dst, const sk_sampling_options_t* sampling);
SK_C_API bool sk_pixmap_erase_color(const sk_pixmap_t* cpixmap, sk_color_t color, const sk_irect_t* subset);
SK_C_API bool sk_pixmap_erase_color4f(const sk_pixmap_t* cpixmap, const sk_color4f_t* color, const sk_irect_t* subset);

// Sk*Encoder

SK_C_API bool sk_webpencoder_encode(sk_wstream_t* dst, const sk_pixmap_t* src, const sk_webpencoder_options_t* options);
SK_C_API bool sk_jpegencoder_encode(sk_wstream_t* dst, const sk_pixmap_t* src, const sk_jpegencoder_options_t* options);
SK_C_API bool sk_pngencoder_encode(sk_wstream_t* dst, const sk_pixmap_t* src, const sk_pngencoder_options_t* options);

// SkSwizzle

SK_C_API void sk_swizzle_swap_rb(uint32_t* dest, const uint32_t* src, int count);

// SkColor

SK_C_API sk_color_t sk_color_unpremultiply(const sk_pmcolor_t pmcolor);
SK_C_API sk_pmcolor_t sk_color_premultiply(const sk_color_t color);
SK_C_API void sk_color_unpremultiply_array(const sk_pmcolor_t* pmcolors, int size, sk_color_t* colors);
SK_C_API void sk_color_premultiply_array(const sk_color_t* colors, int size, sk_pmcolor_t* pmcolors);
SK_C_API void sk_color_get_bit_shift(int* a, int* r, int* g, int* b);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_region.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2016 Bluebeam Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_region_DEFINED
#define sk_region_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

// sk_region_t

SK_C_API sk_region_t* sk_region_new(void);
SK_C_API void sk_region_delete(sk_region_t* r);
SK_C_API bool sk_region_is_empty(const sk_region_t* r);
SK_C_API bool sk_region_is_rect(const sk_region_t* r);
SK_C_API bool sk_region_is_complex(const sk_region_t* r);
SK_C_API void sk_region_get_bounds(const sk_region_t* r, sk_irect_t* rect);
SK_C_API bool sk_region_get_boundary_path(const sk_region_t* r, sk_path_t* path);
SK_C_API bool sk_region_set_empty(sk_region_t* r);
SK_C_API bool sk_region_set_rect(sk_region_t* r, const sk_irect_t* rect);
SK_C_API bool sk_region_set_rects(sk_region_t* r, const sk_irect_t* rects, int count);
SK_C_API bool sk_region_set_region(sk_region_t* r, const sk_region_t* region);
SK_C_API bool sk_region_set_path(sk_region_t* r, const sk_path_t* t, const sk_region_t* clip);
SK_C_API bool sk_region_intersects_rect(const sk_region_t* r, const sk_irect_t* rect);
SK_C_API bool sk_region_intersects(const sk_region_t* r, const sk_region_t* src);
SK_C_API bool sk_region_contains_point(const sk_region_t* r, int x, int y);
SK_C_API bool sk_region_contains_rect(const sk_region_t* r, const sk_irect_t* rect);
SK_C_API bool sk_region_contains(const sk_region_t* r, const sk_region_t* region);
SK_C_API bool sk_region_quick_contains(const sk_region_t* r, const sk_irect_t* rect);
SK_C_API bool sk_region_quick_reject_rect(const sk_region_t* r, const sk_irect_t* rect);
SK_C_API bool sk_region_quick_reject(const sk_region_t* r, const sk_region_t* region);
SK_C_API void sk_region_translate(sk_region_t* r, int x, int y);
SK_C_API bool sk_region_op_rect(sk_region_t* r, const sk_irect_t* rect, sk_region_op_t op);
SK_C_API bool sk_region_op(sk_region_t* r, const sk_region_t* region, sk_region_op_t op);

// sk_region_iterator_t

SK_C_API sk_region_iterator_t* sk_region_iterator_new(const sk_region_t* region);
SK_C_API void sk_region_iterator_delete(sk_region_iterator_t* iter);
SK_C_API bool sk_region_iterator_rewind(sk_region_iterator_t* iter);
SK_C_API bool sk_region_iterator_done(const sk_region_iterator_t* iter);
SK_C_API void sk_region_iterator_next(sk_region_iterator_t* iter);
SK_C_API void sk_region_iterator_rect(const sk_region_iterator_t* iter, sk_irect_t* rect);

// sk_region_cliperator_t

SK_C_API sk_region_cliperator_t* sk_region_cliperator_new(const sk_region_t* region, const sk_irect_t* clip);
SK_C_API void sk_region_cliperator_delete(sk_region_cliperator_t* iter);
SK_C_API bool sk_region_cliperator_done(sk_region_cliperator_t* iter);
SK_C_API void sk_region_cliperator_next(sk_region_cliperator_t* iter);
SK_C_API void sk_region_cliperator_rect(const sk_region_cliperator_t* iter, sk_irect_t* rect);

// sk_region_spanerator_t

SK_C_API sk_region_spanerator_t* sk_region_spanerator_new(const sk_region_t* region, int y, int left, int right);
SK_C_API void sk_region_spanerator_delete(sk_region_spanerator_t* iter);
SK_C_API bool sk_region_spanerator_next(sk_region_spanerator_t* iter, int* left, int* right);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_rrect.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2016 Xamarin Inc.
 * Copyright 2018 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_rrect_DEFINED
#define sk_rrect_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API sk_rrect_t* sk_rrect_new(void);
SK_C_API sk_rrect_t* sk_rrect_new_copy(const sk_rrect_t* rrect);
SK_C_API void sk_rrect_delete(const sk_rrect_t* rrect);
SK_C_API sk_rrect_type_t sk_rrect_get_type(const sk_rrect_t* rrect);
SK_C_API void sk_rrect_get_rect(const sk_rrect_t* rrect, sk_rect_t* rect);
SK_C_API void sk_rrect_get_radii(const sk_rrect_t* rrect, sk_rrect_corner_t corner, sk_vector_t* radii);
SK_C_API float sk_rrect_get_width(const sk_rrect_t* rrect);
SK_C_API float sk_rrect_get_height(const sk_rrect_t* rrect);
SK_C_API void sk_rrect_set_empty(sk_rrect_t* rrect);
SK_C_API void sk_rrect_set_rect(sk_rrect_t* rrect, const sk_rect_t* rect);
SK_C_API void sk_rrect_set_oval(sk_rrect_t* rrect, const sk_rect_t* rect);
SK_C_API void sk_rrect_set_rect_xy(sk_rrect_t* rrect, const sk_rect_t* rect, float xRad, float yRad);
SK_C_API void sk_rrect_set_nine_patch(sk_rrect_t* rrect, const sk_rect_t* rect, float leftRad, float topRad, float rightRad, float bottomRad);
SK_C_API void sk_rrect_set_rect_radii(sk_rrect_t* rrect, const sk_rect_t* rect, const sk_vector_t* radii);
SK_C_API void sk_rrect_inset(sk_rrect_t* rrect, float dx, float dy);
SK_C_API void sk_rrect_outset(sk_rrect_t* rrect, float dx, float dy);
SK_C_API void sk_rrect_offset(sk_rrect_t* rrect, float dx, float dy);
SK_C_API bool sk_rrect_contains(const sk_rrect_t* rrect, const sk_rect_t* rect);
SK_C_API bool sk_rrect_is_valid(const sk_rrect_t* rrect);
SK_C_API bool sk_rrect_transform(sk_rrect_t* rrect, const sk_matrix_t* matrix, sk_rrect_t* dest);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_runtimeeffect.h
/*
 * Copyright 2020 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_runtimeeffect_DEFINED
#define sk_runtimeeffect_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API sk_runtimeeffect_t* sk_runtimeeffect_make_for_color_filter(sk_string_t* sksl, sk_string_t* error);
SK_C_API sk_runtimeeffect_t* sk_runtimeeffect_make_for_shader(sk_string_t* sksl, sk_string_t* error);
SK_C_API sk_runtimeeffect_t* sk_runtimeeffect_make_for_blender(sk_string_t* sksl, sk_string_t* error);
SK_C_API void sk_runtimeeffect_unref(sk_runtimeeffect_t* effect);
SK_C_API sk_shader_t* sk_runtimeeffect_make_shader(sk_runtimeeffect_t* effect, sk_data_t* uniforms, sk_flattenable_t** children, size_t childCount, const sk_matrix_t* localMatrix);
SK_C_API sk_colorfilter_t* sk_runtimeeffect_make_color_filter(sk_runtimeeffect_t* effect, sk_data_t* uniforms, sk_flattenable_t** children, size_t childCount);
SK_C_API sk_blender_t* sk_runtimeeffect_make_blender(sk_runtimeeffect_t* effect, sk_data_t* uniforms, sk_flattenable_t** children, size_t childCount);
SK_C_API size_t sk_runtimeeffect_get_uniform_byte_size(const sk_runtimeeffect_t* effect);

SK_C_API size_t sk_runtimeeffect_get_uniforms_size(const sk_runtimeeffect_t* effect);
SK_C_API void sk_runtimeeffect_get_uniform_name(const sk_runtimeeffect_t* effect, int index, sk_string_t* name);
SK_C_API void sk_runtimeeffect_get_uniform_from_index(const sk_runtimeeffect_t* effect, int index, sk_runtimeeffect_uniform_t* cuniform);
SK_C_API void sk_runtimeeffect_get_uniform_from_name(const sk_runtimeeffect_t* effect, const char* name, size_t len, sk_runtimeeffect_uniform_t* cuniform);

SK_C_API size_t sk_runtimeeffect_get_children_size(const sk_runtimeeffect_t* effect);
SK_C_API void sk_runtimeeffect_get_child_name(const sk_runtimeeffect_t* effect, int index, sk_string_t* name);
SK_C_API void sk_runtimeeffect_get_child_from_index(const sk_runtimeeffect_t* effect, int index, sk_runtimeeffect_child_t* cchild);
SK_C_API void sk_runtimeeffect_get_child_from_name(const sk_runtimeeffect_t* effect, const char* name, size_t len, sk_runtimeeffect_child_t* cchild);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_shader.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_shader_DEFINED
#define sk_shader_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

// SkShader

SK_C_API void sk_shader_ref(sk_shader_t* shader);
SK_C_API void sk_shader_unref(sk_shader_t* shader);
SK_C_API sk_shader_t* sk_shader_with_local_matrix(const sk_shader_t* shader, const sk_matrix_t* localMatrix);
SK_C_API sk_shader_t* sk_shader_with_color_filter(const sk_shader_t* shader, const sk_colorfilter_t* filter);

// SkShaders

SK_C_API sk_shader_t* sk_shader_new_empty(void);
SK_C_API sk_shader_t* sk_shader_new_color(sk_color_t color);
SK_C_API sk_shader_t* sk_shader_new_color4f(const sk_color4f_t* color, const sk_colorspace_t* colorspace);
SK_C_API sk_shader_t* sk_shader_new_blend(sk_blendmode_t mode, const sk_shader_t* dst, const sk_shader_t* src);
SK_C_API sk_shader_t* sk_shader_new_blender(sk_blender_t* blender, const sk_shader_t* dst, const sk_shader_t* src);

// SkGradientShader

SK_C_API sk_shader_t* sk_shader_new_linear_gradient(const sk_point_t points[2], const sk_color_t colors[], const float colorPos[], int colorCount, sk_shader_tilemode_t tileMode, const sk_matrix_t* localMatrix);
SK_C_API sk_shader_t* sk_shader_new_linear_gradient_color4f(const sk_point_t points[2], const sk_color4f_t* colors, const sk_colorspace_t* colorspace, const float colorPos[], int colorCount, sk_shader_tilemode_t tileMode, const sk_matrix_t* localMatrix);
SK_C_API sk_shader_t* sk_shader_new_radial_gradient(const sk_point_t* center, float radius, const sk_color_t colors[], const float colorPos[], int colorCount, sk_shader_tilemode_t tileMode, const sk_matrix_t* localMatrix);
SK_C_API sk_shader_t* sk_shader_new_radial_gradient_color4f(const sk_point_t* center, float radius, const sk_color4f_t* colors, const sk_colorspace_t* colorspace, const float colorPos[], int colorCount, sk_shader_tilemode_t tileMode, const sk_matrix_t* localMatrix);
SK_C_API sk_shader_t* sk_shader_new_sweep_gradient(const sk_point_t* center, const sk_color_t colors[], const float colorPos[], int colorCount, sk_shader_tilemode_t tileMode, float startAngle, float endAngle, const sk_matrix_t* localMatrix);
SK_C_API sk_shader_t* sk_shader_new_sweep_gradient_color4f(const sk_point_t* center, const sk_color4f_t* colors, const sk_colorspace_t* colorspace, const float colorPos[], int colorCount, sk_shader_tilemode_t tileMode, float startAngle, float endAngle, const sk_matrix_t* localMatrix);
SK_C_API sk_shader_t* sk_shader_new_two_point_conical_gradient(const sk_point_t* start, float startRadius, const sk_point_t* end, float endRadius, const sk_color_t colors[], const float colorPos[], int colorCount, sk_shader_tilemode_t tileMode, const sk_matrix_t* localMatrix);
SK_C_API sk_shader_t* sk_shader_new_two_point_conical_gradient_color4f(const sk_point_t* start, float startRadius, const sk_point_t* end, float endRadius, const sk_color4f_t* colors, const sk_colorspace_t* colorspace, const float colorPos[], int colorCount, sk_shader_tilemode_t tileMode, const sk_matrix_t* localMatrix);

// SkPerlinNoiseShader

SK_C_API sk_shader_t* sk_shader_new_perlin_noise_fractal_noise(float baseFrequencyX, float baseFrequencyY, int numOctaves, float seed, const sk_isize_t* tileSize);
SK_C_API sk_shader_t* sk_shader_new_perlin_noise_turbulence(float baseFrequencyX, float baseFrequencyY, int numOctaves, float seed, const sk_isize_t* tileSize);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_stream.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_stream_DEFINED
#define sk_stream_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void sk_stream_asset_destroy(sk_stream_asset_t* cstream);

////////////////////////////////////////////////////////////////////////////////

SK_C_API sk_stream_filestream_t* sk_filestream_new(const char* path);
SK_C_API void sk_filestream_destroy(sk_stream_filestream_t* cstream);
SK_C_API bool sk_filestream_is_valid(sk_stream_filestream_t* cstream);

////////////////////////////////////////////////////////////////////////////////

SK_C_API sk_stream_memorystream_t* sk_memorystream_new(void);
SK_C_API sk_stream_memorystream_t* sk_memorystream_new_with_length(size_t length);
SK_C_API sk_stream_memorystream_t* sk_memorystream_new_with_data(const void* data, size_t length, bool copyData);
SK_C_API sk_stream_memorystream_t* sk_memorystream_new_with_skdata(sk_data_t* data);
SK_C_API void sk_memorystream_set_memory(sk_stream_memorystream_t* cmemorystream, const void* data, size_t length, bool copyData);
SK_C_API void sk_memorystream_destroy(sk_stream_memorystream_t* cstream);

////////////////////////////////////////////////////////////////////////////////

SK_C_API size_t sk_stream_read(sk_stream_t* cstream, void* buffer, size_t size);
SK_C_API size_t sk_stream_peek(sk_stream_t* cstream, void* buffer, size_t size);
SK_C_API size_t sk_stream_skip(sk_stream_t* cstream, size_t size);
SK_C_API bool sk_stream_is_at_end(sk_stream_t* cstream);
SK_C_API bool sk_stream_read_s8(sk_stream_t* cstream, int8_t* buffer);
SK_C_API bool sk_stream_read_s16(sk_stream_t* cstream, int16_t* buffer);
SK_C_API bool sk_stream_read_s32(sk_stream_t* cstream, int32_t* buffer);
SK_C_API bool sk_stream_read_u8(sk_stream_t* cstream, uint8_t* buffer);
SK_C_API bool sk_stream_read_u16(sk_stream_t* cstream, uint16_t* buffer);
SK_C_API bool sk_stream_read_u32(sk_stream_t* cstream, uint32_t* buffer);
SK_C_API bool sk_stream_read_bool(sk_stream_t* cstream, bool* buffer);
SK_C_API bool sk_stream_rewind(sk_stream_t* cstream);
SK_C_API bool sk_stream_has_position(sk_stream_t* cstream);
SK_C_API size_t sk_stream_get_position(sk_stream_t* cstream);
SK_C_API bool sk_stream_seek(sk_stream_t* cstream, size_t position);
SK_C_API bool sk_stream_move(sk_stream_t* cstream, long offset);
SK_C_API bool sk_stream_has_length(sk_stream_t* cstream);
SK_C_API size_t sk_stream_get_length(sk_stream_t* cstream);
SK_C_API const void* sk_stream_get_memory_base(sk_stream_t* cstream);
SK_C_API sk_stream_t* sk_stream_fork(sk_stream_t* cstream);
SK_C_API sk_stream_t* sk_stream_duplicate(sk_stream_t* cstream);
SK_C_API void sk_stream_destroy(sk_stream_t* cstream);

////////////////////////////////////////////////////////////////////////////////

SK_C_API sk_wstream_filestream_t* sk_filewstream_new(const char* path);
SK_C_API void sk_filewstream_destroy(sk_wstream_filestream_t* cstream);
SK_C_API bool sk_filewstream_is_valid(sk_wstream_filestream_t* cstream);

SK_C_API sk_wstream_dynamicmemorystream_t* sk_dynamicmemorywstream_new(void);
SK_C_API sk_stream_asset_t* sk_dynamicmemorywstream_detach_as_stream(sk_wstream_dynamicmemorystream_t* cstream);
SK_C_API sk_data_t* sk_dynamicmemorywstream_detach_as_data(sk_wstream_dynamicmemorystream_t* cstream);
SK_C_API void sk_dynamicmemorywstream_copy_to(sk_wstream_dynamicmemorystream_t* cstream, void* data);
SK_C_API bool sk_dynamicmemorywstream_write_to_stream(sk_wstream_dynamicmemorystream_t* cstream, sk_wstream_t* dst);
SK_C_API void sk_dynamicmemorywstream_destroy(sk_wstream_dynamicmemorystream_t* cstream);

////////////////////////////////////////////////////////////////////////////////

SK_C_API bool sk_wstream_write(sk_wstream_t* cstream, const void* buffer, size_t size);
SK_C_API bool sk_wstream_newline(sk_wstream_t* cstream);
SK_C_API void sk_wstream_flush(sk_wstream_t* cstream);
SK_C_API size_t sk_wstream_bytes_written(sk_wstream_t* cstream);
SK_C_API bool sk_wstream_write_8(sk_wstream_t* cstream, uint8_t value);
SK_C_API bool sk_wstream_write_16(sk_wstream_t* cstream, uint16_t value);
SK_C_API bool sk_wstream_write_32(sk_wstream_t* cstream, uint32_t value);
SK_C_API bool sk_wstream_write_text(sk_wstream_t* cstream, const char* value);
SK_C_API bool sk_wstream_write_dec_as_text(sk_wstream_t* cstream, int32_t value);
SK_C_API bool sk_wstream_write_bigdec_as_text(sk_wstream_t* cstream, int64_t value, int minDigits);
SK_C_API bool sk_wstream_write_hex_as_text(sk_wstream_t* cstream, uint32_t value, int minDigits);
SK_C_API bool sk_wstream_write_scalar_as_text(sk_wstream_t* cstream, float value);
SK_C_API bool sk_wstream_write_bool(sk_wstream_t* cstream, bool value);
SK_C_API bool sk_wstream_write_scalar(sk_wstream_t* cstream, float value);
SK_C_API bool sk_wstream_write_packed_uint(sk_wstream_t* cstream, size_t value);
SK_C_API bool sk_wstream_write_stream(sk_wstream_t* cstream, sk_stream_t* input, size_t length);
SK_C_API int sk_wstream_get_size_of_packed_uint(size_t value);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_string.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_string_DEFINED
#define sk_string_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API sk_string_t* sk_string_new_empty(void);
SK_C_API sk_string_t* sk_string_new_with_copy(const char* src, size_t length);
SK_C_API void sk_string_destructor(const sk_string_t*);
SK_C_API size_t sk_string_get_size(const sk_string_t*);
SK_C_API const char* sk_string_get_c_str(const sk_string_t*);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_surface.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_surface_DEFINED
#define sk_surface_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

// surface

SK_C_API sk_surface_t* sk_surface_new_null(int width, int height);
SK_C_API sk_surface_t* sk_surface_new_raster(const sk_imageinfo_t*, size_t rowBytes, const sk_surfaceprops_t*);
SK_C_API sk_surface_t* sk_surface_new_raster_direct(const sk_imageinfo_t*, void* pixels, size_t rowBytes, const sk_surface_raster_release_proc releaseProc, void* context, const sk_surfaceprops_t* props);

SK_C_API sk_surface_t* sk_surface_new_backend_texture(gr_recording_context_t* context, const gr_backendtexture_t* texture, gr_surfaceorigin_t origin, int samples, sk_colortype_t colorType, sk_colorspace_t* colorspace, const sk_surfaceprops_t* props);
SK_C_API sk_surface_t* sk_surface_new_backend_render_target(gr_recording_context_t* context, const gr_backendrendertarget_t* target, gr_surfaceorigin_t origin, sk_colortype_t colorType, sk_colorspace_t* colorspace, const sk_surfaceprops_t* props);
SK_C_API sk_surface_t* sk_surface_new_render_target(gr_recording_context_t* context, bool budgeted, const sk_imageinfo_t* cinfo, int sampleCount, gr_surfaceorigin_t origin, const sk_surfaceprops_t* props, bool shouldCreateWithMips);

SK_C_API sk_surface_t* sk_surface_new_metal_layer(gr_recording_context_t* context, gr_mtl_handle_t layer, gr_surfaceorigin_t origin, int sampleCount, sk_colortype_t colorType, sk_colorspace_t* colorspace, const sk_surfaceprops_t* props, gr_mtl_handle_t* drawable);
SK_C_API sk_surface_t* sk_surface_new_metal_view(gr_recording_context_t* context, gr_mtl_handle_t mtkView, gr_surfaceorigin_t origin, int sampleCount, sk_colortype_t colorType, sk_colorspace_t* colorspace, const sk_surfaceprops_t* props);

SK_C_API void sk_surface_unref(sk_surface_t*);
SK_C_API sk_canvas_t* sk_surface_get_canvas(sk_surface_t*);
SK_C_API sk_image_t* sk_surface_new_image_snapshot(sk_surface_t*);
SK_C_API sk_image_t* sk_surface_new_image_snapshot_with_crop(sk_surface_t* surface, const sk_irect_t* bounds);
SK_C_API void sk_surface_draw(sk_surface_t* surface, sk_canvas_t* canvas, float x, float y, const sk_paint_t* paint);
SK_C_API bool sk_surface_peek_pixels(sk_surface_t* surface, sk_pixmap_t* pixmap);
SK_C_API bool sk_surface_read_pixels(sk_surface_t* surface, sk_imageinfo_t* dstInfo, void* dstPixels, size_t dstRowBytes, int srcX, int srcY);
SK_C_API const sk_surfaceprops_t* sk_surface_get_props(sk_surface_t* surface);
SK_C_API void sk_surface_flush(sk_surface_t* surface);
SK_C_API void sk_surface_flush_and_submit(sk_surface_t* surface, bool syncCpu);
SK_C_API gr_recording_context_t* sk_surface_get_recording_context(sk_surface_t* surface);

// surface props

SK_C_API sk_surfaceprops_t* sk_surfaceprops_new(uint32_t flags, sk_pixelgeometry_t geometry);
SK_C_API void sk_surfaceprops_delete(sk_surfaceprops_t* props);
SK_C_API uint32_t sk_surfaceprops_get_flags(sk_surfaceprops_t* props);
SK_C_API sk_pixelgeometry_t sk_surfaceprops_get_pixel_geometry(sk_surfaceprops_t* props);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_svg.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_svg_DEFINED
#define sk_svg_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API sk_canvas_t* sk_svgcanvas_create_with_stream(const sk_rect_t* bounds, sk_wstream_t* stream);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_textblob.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_textblob_DEFINED
#define sk_textblob_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

// sk_textblob_t

SK_C_API void sk_textblob_ref(const sk_textblob_t* blob);
SK_C_API void sk_textblob_unref(const sk_textblob_t* blob);
SK_C_API uint32_t sk_textblob_get_unique_id(const sk_textblob_t* blob);
SK_C_API void sk_textblob_get_bounds(const sk_textblob_t* blob, sk_rect_t* bounds);
SK_C_API int sk_textblob_get_intercepts(const sk_textblob_t* blob, const float bounds[2], float intervals[], const sk_paint_t* paint);

// sk_textblob_builder_t

SK_C_API sk_textblob_builder_t* sk_textblob_builder_new(void);
SK_C_API void sk_textblob_builder_delete(sk_textblob_builder_t* builder);
SK_C_API sk_textblob_t* sk_textblob_builder_make(sk_textblob_builder_t* builder);
SK_C_API void sk_textblob_builder_alloc_run(sk_textblob_builder_t* builder, const sk_font_t* font, int count, float x, float y, const sk_rect_t* bounds, sk_textblob_builder_runbuffer_t* runbuffer);
SK_C_API void sk_textblob_builder_alloc_run_pos_h(sk_textblob_builder_t* builder, const sk_font_t* font, int count, float y, const sk_rect_t* bounds, sk_textblob_builder_runbuffer_t* runbuffer);
SK_C_API void sk_textblob_builder_alloc_run_pos(sk_textblob_builder_t* builder, const sk_font_t* font, int count, const sk_rect_t* bounds, sk_textblob_builder_runbuffer_t* runbuffer);
SK_C_API void sk_textblob_builder_alloc_run_rsxform(sk_textblob_builder_t* builder, const sk_font_t* font, int count, const sk_rect_t* bounds, sk_textblob_builder_runbuffer_t* runbuffer);
SK_C_API void sk_textblob_builder_alloc_run_text(sk_textblob_builder_t* builder, const sk_font_t* font, int count, float x, float y, int textByteCount, const sk_rect_t* bounds, sk_textblob_builder_runbuffer_t* runbuffer);
SK_C_API void sk_textblob_builder_alloc_run_text_pos_h(sk_textblob_builder_t* builder, const sk_font_t* font, int count, float y, int textByteCount, const sk_rect_t* bounds, sk_textblob_builder_runbuffer_t* runbuffer);
SK_C_API void sk_textblob_builder_alloc_run_text_pos(sk_textblob_builder_t* builder, const sk_font_t* font, int count, int textByteCount, const sk_rect_t* bounds, sk_textblob_builder_runbuffer_t* runbuffer);
SK_C_API void sk_textblob_builder_alloc_run_text_rsxform(sk_textblob_builder_t* builder, const sk_font_t* font, int count, int textByteCount, const sk_rect_t* bounds, sk_textblob_builder_runbuffer_t* runbuffer);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_typeface.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_typeface_DEFINED
#define sk_typeface_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

// typeface

SK_C_API void sk_typeface_unref(sk_typeface_t* typeface);
SK_C_API sk_fontstyle_t* sk_typeface_get_fontstyle(const sk_typeface_t* typeface);
SK_C_API int sk_typeface_get_font_weight(const sk_typeface_t* typeface);
SK_C_API int sk_typeface_get_font_width(const sk_typeface_t* typeface);
SK_C_API sk_font_style_slant_t sk_typeface_get_font_slant(const sk_typeface_t* typeface);
SK_C_API bool sk_typeface_is_fixed_pitch(const sk_typeface_t* typeface);
SK_C_API sk_typeface_t* sk_typeface_create_default(void);
SK_C_API sk_typeface_t* sk_typeface_ref_default(void);
SK_C_API sk_typeface_t* sk_typeface_create_from_name(const char* familyName, const sk_fontstyle_t* style);
SK_C_API sk_typeface_t* sk_typeface_create_from_file(const char* path, int index);
SK_C_API sk_typeface_t* sk_typeface_create_from_stream(sk_stream_asset_t* stream, int index);
SK_C_API sk_typeface_t* sk_typeface_create_from_data(sk_data_t* data, int index);
SK_C_API void sk_typeface_unichars_to_glyphs(const sk_typeface_t* typeface, const int32_t unichars[], int count, uint16_t glyphs[]);
SK_C_API uint16_t sk_typeface_unichar_to_glyph(const sk_typeface_t* typeface, const int32_t unichar);
SK_C_API int sk_typeface_count_glyphs(const sk_typeface_t* typeface);
SK_C_API int sk_typeface_count_tables(const sk_typeface_t* typeface);
SK_C_API int sk_typeface_get_table_tags(const sk_typeface_t* typeface, sk_font_table_tag_t tags[]);
SK_C_API size_t sk_typeface_get_table_size(const sk_typeface_t* typeface, sk_font_table_tag_t tag);
SK_C_API size_t sk_typeface_get_table_data(const sk_typeface_t* typeface, sk_font_table_tag_t tag, size_t offset, size_t length, void* data);
SK_C_API sk_data_t* sk_typeface_copy_table_data(const sk_typeface_t* typeface, sk_font_table_tag_t tag);
SK_C_API int sk_typeface_get_units_per_em(const sk_typeface_t* typeface);
SK_C_API bool sk_typeface_get_kerning_pair_adjustments(const sk_typeface_t* typeface, const uint16_t glyphs[], int count, int32_t adjustments[]);
// TODO: createFamilyNameIterator
SK_C_API sk_string_t* sk_typeface_get_family_name(const sk_typeface_t* typeface);
SK_C_API sk_stream_asset_t* sk_typeface_open_stream(const sk_typeface_t* typeface, int* ttcIndex);


// font manager

SK_C_API sk_fontmgr_t* sk_fontmgr_create_default(void);
SK_C_API sk_fontmgr_t* sk_fontmgr_ref_default(void);
SK_C_API void sk_fontmgr_unref(sk_fontmgr_t*);
SK_C_API int sk_fontmgr_count_families(sk_fontmgr_t*);
SK_C_API void sk_fontmgr_get_family_name(sk_fontmgr_t*, int index, sk_string_t* familyName);
SK_C_API sk_fontstyleset_t* sk_fontmgr_create_styleset(sk_fontmgr_t*, int index);
SK_C_API sk_fontstyleset_t* sk_fontmgr_match_family(sk_fontmgr_t*, const char* familyName);
SK_C_API sk_typeface_t* sk_fontmgr_match_family_style(sk_fontmgr_t*, const char* familyName, sk_fontstyle_t* style);
SK_C_API sk_typeface_t* sk_fontmgr_match_family_style_character(sk_fontmgr_t*, const char* familyName, sk_fontstyle_t* style, const char** bcp47, int bcp47Count, int32_t character);
SK_C_API sk_typeface_t* sk_fontmgr_create_from_data(sk_fontmgr_t*, sk_data_t* data, int index);
SK_C_API sk_typeface_t* sk_fontmgr_create_from_stream(sk_fontmgr_t*, sk_stream_asset_t* stream, int index);
SK_C_API sk_typeface_t* sk_fontmgr_create_from_file(sk_fontmgr_t*, const char* path, int index);

// font style

SK_C_API sk_fontstyle_t* sk_fontstyle_new(int weight, int width, sk_font_style_slant_t slant);
SK_C_API void sk_fontstyle_delete(sk_fontstyle_t* fs);
SK_C_API int sk_fontstyle_get_weight(const sk_fontstyle_t* fs);
SK_C_API int sk_fontstyle_get_width(const sk_fontstyle_t* fs);
SK_C_API sk_font_style_slant_t sk_fontstyle_get_slant(const sk_fontstyle_t* fs);

// font style set

SK_C_API sk_fontstyleset_t* sk_fontstyleset_create_empty(void);
SK_C_API void sk_fontstyleset_unref(sk_fontstyleset_t* fss);
SK_C_API int sk_fontstyleset_get_count(sk_fontstyleset_t* fss);
SK_C_API void sk_fontstyleset_get_style(sk_fontstyleset_t* fss, int index, sk_fontstyle_t* fs, sk_string_t* style);
SK_C_API sk_typeface_t* sk_fontstyleset_create_typeface(sk_fontstyleset_t* fss, int index);
SK_C_API sk_typeface_t* sk_fontstyleset_match_style(sk_fontstyleset_t* fss, sk_fontstyle_t* style);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sk_vertices.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sk_vertices_DEFINED
#define sk_vertices_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void sk_vertices_unref(sk_vertices_t* cvertices);
SK_C_API void sk_vertices_ref(sk_vertices_t* cvertices);
SK_C_API sk_vertices_t* sk_vertices_make_copy(sk_vertices_vertex_mode_t vmode, int vertexCount, const sk_point_t* positions, const sk_point_t* texs, const sk_color_t* colors, int indexCount, const uint16_t* indices);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\skottie_animation.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef skottie_DEFINED
#define skottie_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

/*
 * skottie::Animation
 */

SK_C_API skottie_animation_t* skottie_animation_make_from_string(const char* data, size_t length);
SK_C_API skottie_animation_t* skottie_animation_make_from_data(const char* data, size_t length);
SK_C_API skottie_animation_t* skottie_animation_make_from_stream(sk_stream_t* stream);
SK_C_API skottie_animation_t* skottie_animation_make_from_file(const char* path);

SK_C_API void skottie_animation_ref(skottie_animation_t* instance);
SK_C_API void skottie_animation_unref(skottie_animation_t* instance);

SK_C_API void skottie_animation_delete(skottie_animation_t *instance);

SK_C_API void skottie_animation_render(skottie_animation_t *instance, sk_canvas_t *canvas, sk_rect_t *dst);
SK_C_API void skottie_animation_render_with_flags(skottie_animation_t *instance, sk_canvas_t *canvas, sk_rect_t *dst, skottie_animation_renderflags_t flags);

SK_C_API void skottie_animation_seek(skottie_animation_t *instance, float t, sksg_invalidation_controller_t *ic);
SK_C_API void skottie_animation_seek_frame(skottie_animation_t *instance, float t, sksg_invalidation_controller_t *ic);
SK_C_API void skottie_animation_seek_frame_time(skottie_animation_t *instance, float t, sksg_invalidation_controller_t *ic);

SK_C_API double skottie_animation_get_duration(skottie_animation_t *instance);
SK_C_API double skottie_animation_get_fps(skottie_animation_t *instance);
SK_C_API double skottie_animation_get_in_point(skottie_animation_t *instance);
SK_C_API double skottie_animation_get_out_point(skottie_animation_t *instance);

SK_C_API void skottie_animation_get_version(skottie_animation_t *instance, sk_string_t* version);
SK_C_API void skottie_animation_get_size(skottie_animation_t *instance, sk_size_t* size);


/*
 * skottie::Animation::Builder
 */

SK_C_API skottie_animation_builder_t* skottie_animation_builder_new(skottie_animation_builder_flags_t flags);

SK_C_API void skottie_animation_builder_delete(skottie_animation_builder_t *instance);

SK_C_API void skottie_animation_builder_get_stats(skottie_animation_builder_t* instance, skottie_animation_builder_stats_t* stats);
SK_C_API void skottie_animation_builder_set_resource_provider(skottie_animation_builder_t* instance, skottie_resource_provider_t* resourceProvider);
SK_C_API void skottie_animation_builder_set_font_manager(skottie_animation_builder_t* instance, sk_fontmgr_t* fontManager);

SK_C_API skottie_animation_t* skottie_animation_builder_make_from_stream(skottie_animation_builder_t* instance, sk_stream_t* stream);
SK_C_API skottie_animation_t* skottie_animation_builder_make_from_file(skottie_animation_builder_t* instance, const char* path);
SK_C_API skottie_animation_t* skottie_animation_builder_make_from_string(skottie_animation_builder_t* instance, const char* data, size_t length);
SK_C_API skottie_animation_t* skottie_animation_builder_make_from_data(skottie_animation_builder_t* instance, const char* data, size_t length);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\skresources_resource_provider.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef skresources_resource_provider_DEFINED
#define skresources_resource_provider_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API void skresources_resource_provider_ref(skresources_resource_provider_t* instance);
SK_C_API void skresources_resource_provider_unref(skresources_resource_provider_t* instance);
SK_C_API void skresources_resource_provider_delete(skresources_resource_provider_t *instance);

SK_C_API sk_data_t* skresources_resource_provider_load(skresources_resource_provider_t *instance, const char* path, const char* name);
SK_C_API skresources_image_asset_t* skresources_resource_provider_load_image_asset(skresources_resource_provider_t *instance, const char* path, const char* name, const char* id);
SK_C_API skresources_external_track_asset_t* skresources_resource_provider_load_audio_asset(skresources_resource_provider_t *instance, const char* path, const char* name, const char* id);
SK_C_API sk_typeface_t* skresources_resource_provider_load_typeface(skresources_resource_provider_t *instance, const char* name, const char* url);

SK_C_API skresources_resource_provider_t* skresources_file_resource_provider_make(sk_string_t* base_dir, bool predecode);
SK_C_API skresources_resource_provider_t* skresources_caching_resource_provider_proxy_make(skresources_resource_provider_t* rp);
SK_C_API skresources_resource_provider_t* skresources_data_uri_resource_provider_proxy_make(skresources_resource_provider_t* rp, bool predecode);

SK_C_PLUS_PLUS_END_GUARD

#endif

//c\sksg_invalidation_controller.h
/*
 * Copyright 2014 Google Inc.
 * Copyright 2015 Xamarin Inc.
 * Copyright 2017 Microsoft Corporation. All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef sksg_invalidationcontroller_DEFINED
#define sksg_invalidationcontroller_DEFINED



SK_C_PLUS_PLUS_BEGIN_GUARD

SK_C_API sksg_invalidation_controller_t* sksg_invalidation_controller_new(void);
SK_C_API void sksg_invalidation_controller_delete(sksg_invalidation_controller_t* instance);

SK_C_API void sksg_invalidation_controller_inval(sksg_invalidation_controller_t* instance, sk_rect_t* rect, sk_matrix_t* matrix);
SK_C_API void sksg_invalidation_controller_get_bounds(sksg_invalidation_controller_t* instance, sk_rect_t* bounds);
SK_C_API void sksg_invalidation_controller_begin(sksg_invalidation_controller_t* instance);
SK_C_API void sksg_invalidation_controller_end(sksg_invalidation_controller_t* instance);
SK_C_API void sksg_invalidation_controller_reset(sksg_invalidation_controller_t* instance);

SK_C_PLUS_PLUS_END_GUARD

#endif



